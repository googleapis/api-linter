// Copyright 2019 Google LLC
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

package aip0134

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestResourceField(t *testing.T) {
	// Set up the testing permutations.
	tests := []struct {
		name              string
		MessageName       string
		ResourceName      string
		ResourceFieldName string
		problems          testutils.Problems
	}{
		{"Valid", "UpdateBookRequest", "Book", "book", nil},
		{"InvalidFieldName", "UpdateBookRequest", "Book", "big_book", testutils.Problems{{Suggestion: "book"}}},
		{"IrrelevantMessage", "ModifyBookRequest", "Book", "big_book", nil},
		{"IrrelevantFieldType", "UpdateBookRequest", "string", "big_book", nil},
		{"IrrelevantFieldMessage", "UpdateBookRequest", "Foo", "big_book", nil},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
				import "google/protobuf/field_mask.proto";
				message {{.MessageName}} {
					{{.ResourceName}} {{.ResourceFieldName}} = 1;
				}
				message {{.ResourceName}} {}
			`, test)
			field := file.GetMessageTypes()[0].GetFields()[0]

			// Run the lint rule, and establish that it returns the correct problems.
			problems := requestResourceField.Lint(file)
			if diff := test.problems.SetDescriptor(field).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}
