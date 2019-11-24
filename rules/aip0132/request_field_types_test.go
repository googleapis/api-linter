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

package aip0132

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
	"github.com/jhump/protoreflect/desc/builder"
)

func TestRequestFieldTypes(t *testing.T) {
	// Set up the testing permutations.
	tests := []struct {
		testName    string
		messageName string
		fieldName   string
		fieldType   *builder.FieldType
		problems    testutils.Problems
	}{
		{"Filter", "ListBooksRequest", "filter", builder.FieldTypeString(), testutils.Problems{}},
		{"FilterInvalid", "ListBooksRequest", "filter", builder.FieldTypeBytes(), testutils.Problems{{Message: "string"}}},
		{"OrderBy", "ListBooksRequest", "order_by", builder.FieldTypeString(), testutils.Problems{}},
		{"OrderByInvalid", "ListBooksRequest", "order_by", builder.FieldTypeBytes(), testutils.Problems{{Message: "string"}}},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Create an appropriate message descriptor.
			message, err := builder.NewMessage(test.messageName).AddField(
				builder.NewField("parent", builder.FieldTypeString()),
			).AddField(
				builder.NewField(test.fieldName, test.fieldType),
			).Build()
			if err != nil {
				t.Fatalf("Could not build GetBookRequest message.")
			}

			// Run the lint rule, and establish that it returns the correct
			// number of problems.
			field := message.GetFields()[1]
			problems := requestFieldTypes.Lint(field)
			if diff := test.problems.SetDescriptor(field).Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
