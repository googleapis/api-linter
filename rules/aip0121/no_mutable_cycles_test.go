// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aip0121

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestNoMutableCycles(t *testing.T) {

	for _, test := range []struct {
		name                                                   string
		BookExtensions, PublisherExtensions, LibraryExtensions string
		problems                                               testutils.Problems
	}{
		{
			"ValidNoCycle",
			`[(google.api.resource_reference).type = "library.googleapis.com/Library"]`,
			`[(google.api.resource_reference).type = "library.googleapis.com/Library"]`,
			"",
			nil,
		},
		{
			"InvalidCycle",
			`[(google.api.resource_reference).type = "library.googleapis.com/Publisher"]`,
			`[(google.api.resource_reference).type = "library.googleapis.com/Book"]`,
			"",
			testutils.Problems{{
				Message: "cycle",
			}},
		},
		{
			"InvalidSelfReferenceCycle",
			"",
			`[(google.api.resource_reference).type = "library.googleapis.com/Publisher"]`,
			"",
			testutils.Problems{{
				Message: "cycle",
			}},
		},
		{
			"InvalidDeepCycle",
			`[(google.api.resource_reference).type = "library.googleapis.com/Publisher"]`,
			`[(google.api.resource_reference).type = "library.googleapis.com/Library"]`,
			`[(google.api.resource_reference).type = "library.googleapis.com/Book"]`,
			testutils.Problems{{
				Message: "cycle",
			}},
		},
		{
			"ValidOutputOnlyCyclicReference",
			`[(google.api.resource_reference).type = "library.googleapis.com/Publisher"]`,
			`[
				(google.api.resource_reference).type = "library.googleapis.com/Book",
				(google.api.field_behavior) = OUTPUT_ONLY
			]`,
			"",
			nil,
		},
		{
			"ValidOutputOnlyDeepCyclicReference",
			`[(google.api.resource_reference).type = "library.googleapis.com/Publisher"]`,
			`[(google.api.resource_reference).type = "library.googleapis.com/Library"]`,
			`[
				(google.api.resource_reference).type = "library.googleapis.com/Book",
				(google.api.field_behavior) = OUTPUT_ONLY
			]`,
			nil,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
			import "google/api/resource.proto";
			import "google/api/field_behavior.proto";
			message Book {
				option (google.api.resource) = {
					type: "library.googleapis.com/Book"
					pattern: "publishers/{publisher}/books/{book}"
				};
				string name = 1;

				string resource = 2 {{.BookExtensions}};
			}

			message Publisher {
				option (google.api.resource) = {
					type: "library.googleapis.com/Publisher"
					pattern: "publishers/{publisher}"
				};
				string name = 1;

				string resource = 2 {{.PublisherExtensions}};
			}

			message Library {
				option (google.api.resource) = {
					type: "library.googleapis.com/Library"
					pattern: "libraries/{library}"
				};
				string name = 1;

				string resource = 3 {{.LibraryExtensions}};
			}
			`, test)

			msg := f.GetMessageTypes()[1]
			field := msg.FindFieldByName("resource")
			// If this rule was run on the entire test file, there would be two
			// findings, one for each resource in the cycle. To simplify that,
			// we just lint one of the offending messages.
			if diff := test.problems.SetDescriptor(field).Diff(noMutableCycles.LintMessage(msg)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
