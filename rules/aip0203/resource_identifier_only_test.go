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

package aip0203

import (
	"testing"

	"github.com/googleapis/api-linter/v2/rules/internal/testutils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestResourceIdentifierOnly(t *testing.T) {
	testCases := []struct {
		name, message                                   string
		ResourceField, NameField, NonResourceExtensions string
		problems                                        testutils.Problems
	}{
		{
			name:          "Valid",
			message:       "Book",
			ResourceField: "string name = 1 [(google.api.field_behavior) = IDENTIFIER];",
			problems:      nil,
		},
		{
			name:          "ValidNameField",
			message:       "Book",
			ResourceField: "string resource_name = 1 [(google.api.field_behavior) = IDENTIFIER];",
			NameField:     "resource_name",
			problems:      nil,
		},
		{
			name:                  "InvalidNonResource",
			message:               "NonResource",
			NonResourceExtensions: "[(google.api.field_behavior) = IDENTIFIER]",
			problems:              testutils.Problems{{Message: "resource's name field"}},
		},
		{
			name:          "InvalidResourceNonIdentifier",
			message:       "Book",
			ResourceField: "string foo = 1 [(google.api.field_behavior) = IDENTIFIER];",
			problems:      testutils.Problems{{Message: "resource's name field"}},
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
				{{.ResourceField}}
			}
			message NonResource {
				string foo = 1 {{.NonResourceExtensions}};
			}
			`, test)
			f := file.Messages().ByName(protoreflect.Name(test.message)).Fields().Get(0)
			problems := resourceIdentifierOnly.Lint(file)
			if diff := test.problems.SetDescriptor(f).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}
