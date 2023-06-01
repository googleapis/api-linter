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

package aip0122

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestResourceNamesAlternatingComponents(t *testing.T) {
	for _, test := range []struct {
		name     string
		Pattern  string
		problems testutils.Problems
	}{
		{"ValidProject", "projects/{project}", testutils.Problems{}},
		{"InvalidProjectNoResource", "projects/project", testutils.Problems{{Message: "should usually alternate"}}},
		{"InvalidProjectNoCollection", "{projects}/{project}", testutils.Problems{{Message: "should usually alternate"}}},
		{"ValidLocation", "projects/{project}/locations/{location}", testutils.Problems{}},
		{"InvalidLocationNoResource", "projects/{project}/locations/location", testutils.Problems{{Message: "should usually alternate"}}},
		{"InvalidLocationNoCollection", "{projects}/project/{locations}/{location}", testutils.Problems{{Message: "should usually alternate"}}},
		{"ValidSingleton", "projects/{project}/locations/{location}/settings", testutils.Problems{}},
		{"InvalidSingleton", "projects/{project}/locations/{location}/{settings}", testutils.Problems{{Message: "should usually alternate"}}},
	} {
		f := testutils.ParseProto3Tmpl(t, `
			import "google/api/resource.proto";
			import "google/api/field_behavior.proto";

			message Book {
				option (google.api.resource) = {
					type: "library.googleapis.com/Book"
					pattern: "{{.Pattern}}"
				};
				string project = 1;
				string location = 2;
			}
		`, test)
		message := f.GetMessageTypes()[0]
		if diff := test.problems.SetDescriptor(message).Diff(resourceNamesAlternatingComponents.Lint(f)); diff != "" {
			t.Errorf(diff)
		}
	}
}
