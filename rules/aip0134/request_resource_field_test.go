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
	"github.com/jhump/protoreflect/desc"
)

func TestResourceField(t *testing.T) {
	// Set up the testing permutations.
	tests := []struct {
		testName          string
		MessageName       string
		ResourceName      string
		ResourceFieldName string
		descriptor        func(*desc.MessageDescriptor) desc.Descriptor
		problems          testutils.Problems
	}{
		{"Valid", "UpdateBookRequest", "Book", "book", nil, testutils.Problems{}},
		{"ValidTwoWords", "UpdateBigBookRequest", "BigBook", "big_book", nil, testutils.Problems{}},
		{
			"InvalidMismatch",
			"UpdateBookRequest",
			"Foo",
			"foo",
			nil,
			testutils.Problems{{Message: "has no \"Book\""}}},
		{
			"InvalidFieldName",
			"UpdateBookRequest",
			"Book",
			"big_book",
			func(m *desc.MessageDescriptor) desc.Descriptor {
				return m.GetFields()[0]
			},
			testutils.Problems{{Suggestion: "book"}},
		},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
				import "google/protobuf/field_mask.proto";
				message {{.MessageName}} {
					{{.ResourceName}} {{.ResourceFieldName}} = 1;
					google.protobuf.FieldMask update_mask = 2;
				}
				message {{.ResourceName}} {}
			`, test)
			var d desc.Descriptor = file.GetMessageTypes()[0]
			if test.descriptor != nil {
				d = test.descriptor(d.(*desc.MessageDescriptor))
			}

			// Run the lint rule, and establish that it returns the correct problems.
			problems := resourceField.Lint(file.GetMessageTypes()[0])
			if diff := test.problems.SetDescriptor(d).Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
