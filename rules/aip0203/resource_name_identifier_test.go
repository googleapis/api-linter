// Copyright 2020 Google LLC
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

package aip0203

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestResourceNameIdentifier(t *testing.T) {
	testCases := []struct {
		name             string
		Field, NameField string
		problems         testutils.Problems
	}{
		{
			name:     "Valid",
			Field:    "string name = 1 [(google.api.field_behavior) = IDENTIFIER];",
			problems: nil,
		},
		{
			name:      "ValidNameField",
			Field:     "string resource_name = 1 [(google.api.field_behavior) = IDENTIFIER];",
			NameField: "resource_name",
			problems:  nil,
		},
		{
			name:     "SkipMissingNameField",
			Field:    "string other = 1;",
			problems: nil,
		},
		{
			name:     "InvalidNoFieldBehavior",
			Field:    "string name = 1;",
			problems: testutils.Problems{{Message: "field_behavior IDENTIFIER"}},
		},
		{
			name:     "InvalidMissingIdentifier",
			Field:    "string name = 1 [(google.api.field_behavior) = REQUIRED];",
			problems: testutils.Problems{{Message: "field_behavior IDENTIFIER"}},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
			import "google/api/field_behavior.proto";
			import "google/api/resource.proto";
				message Book {
					option (google.api.resource) = {
						type: "library.googleapis.com/Book"
						pattern: "books/{book}"
						name_field: "{{.NameField}}"
					};

					{{.Field}}
				}`, test)
			f := file.Messages().Get(0).Fields().Get(0)
			problems := resourceNameIdentifier.Lint(file)
			if diff := test.problems.SetDescriptor(f).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}
