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

package aip0135

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
	"github.com/jhump/protoreflect/desc/builder"
)

func TestStandardFields(t *testing.T) {
	// Set up the testing permutations.
	tests := []struct {
		testName      string
		messageName   string
		nameFieldName string
		problems      testutils.Problems
	}{
		{"Valid", "DeleteBookRequest", "name", testutils.Problems{}},
		{"InvalidName", "DeleteBookRequest", "id", testutils.Problems{{Message: "name"}}},
		{"Irrelevant", "RemoveBookRequest", "id", testutils.Problems{}},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Create an appropriate message descriptor.
			message, err := builder.NewMessage(test.messageName).AddField(
				builder.NewField(test.nameFieldName, builder.FieldTypeString()),
			).Build()
			if err != nil {
				t.Fatalf("Could not build %s message.", test.messageName)
			}

			// Run the lint rule, and establish that it returns the correct problems.
			problems := standardFields.Lint(message.GetFile())
			if diff := test.problems.SetDescriptor(message).Diff(problems); diff != "" {
				t.Errorf("Problems did not match: %v", diff)
			}
		})
	}
}

func TestStandardFieldsInvalidType(t *testing.T) {
	// Create an appropriate message descriptor.
	message, err := builder.NewMessage("DeleteBookRequest").AddField(
		builder.NewField("name", builder.FieldTypeBytes()),
	).Build()
	if err != nil {
		t.Fatalf("Could not build descriptor.")
	}

	// Run the lint rule, and establish that it returns the correct
	// number of problems.
	wantProblems := testutils.Problems{{
		Descriptor: message.GetFields()[0],
		Message:    "string",
	}}
	gotProblems := standardFields.Lint(message.GetFile())
	if diff := wantProblems.Diff(gotProblems); diff != "" {
		t.Errorf(diff)
	}
}
