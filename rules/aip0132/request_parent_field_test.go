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
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"
)

func TestStandardFields(t *testing.T) {
	// Set up the testing permutations.
	tests := []struct {
		testName      string
		messageName   string
		nameFieldName string
		nameFieldType *builder.FieldType
		problems      testutils.Problems
		problemDesc   func(m *desc.MessageDescriptor) desc.Descriptor
	}{
		{"Valid", "ListBooksRequest", "parent", builder.FieldTypeString(), testutils.Problems{}, nil},
		{"InvalidName", "ListBooksRequest", "publisher", builder.FieldTypeString(), testutils.Problems{{Message: "no `parent` field"}}, nil},
		{
			"InvalidType",
			"ListBooksRequest",
			"parent",
			builder.FieldTypeBytes(),
			testutils.Problems{{Suggestion: "string"}},
			func(m *desc.MessageDescriptor) desc.Descriptor {
				return m.GetFields()[0]
			},
		},
		{"Irrelevant", "EnumerateBooksRequest", "id", builder.FieldTypeString(), testutils.Problems{}, nil},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Create an appropriate message descriptor.
			message, err := builder.NewMessage(test.messageName).AddField(
				builder.NewField(test.nameFieldName, test.nameFieldType),
			).Build()
			if err != nil {
				t.Fatalf("Could not build %s message.", test.messageName)
			}

			// What descriptor is the problem expected to be attached to?
			var problemDesc desc.Descriptor = message
			if test.problemDesc != nil {
				problemDesc = test.problemDesc(message)
			}

			// Run the lint rule, and establish that it returns the correct problems.
			// number of problems.
			problems := standardFields.Lint(message)
			if diff := test.problems.SetDescriptor(problemDesc).Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
