// Copyright 2021 Google LLC
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

package aip0123

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestResourceReferenceType(t *testing.T) {
	annotation := ` [(google.api.resource_reference) = {
		type: "library.googleapis.com/Author"
	}]`
	for _, test := range []struct {
		name       string
		Type       string
		Annotation string
		problems   testutils.Problems
	}{
		{"Valid", "string", annotation, nil},
		{"Invalid", "Author", annotation, testutils.Problems{{Suggestion: ""}}},
		{"Irrelevant", "Author", "", nil},
	} {
		t.Run(test.name, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
			  import "google/api/resource.proto";
				message Book {
					{{.Type}} author = 1{{.Annotation}};
				}
				message Author {}  // For the failure case, to make compilation work.
			`, test)
			field := file.GetMessageTypes()[0].GetFields()[0]
			got := resourceReferenceType.Lint(file)
			if diff := test.problems.SetDescriptor(field).Diff(got); diff != "" {
				t.Error(diff)
			}
		})
	}
}
