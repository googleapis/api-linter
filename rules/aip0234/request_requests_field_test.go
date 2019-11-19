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

package aip0234

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
	"github.com/jhump/protoreflect/desc"
)

func TestRequestRequestsField(t *testing.T) {
	// Set up the testing permutations.
	tests := []struct {
		testName string
		Field    string
		problems testutils.Problems
	}{
		{
			testName: "Valid",
			Field:    "repeated UpdateBookRequest requests",
			problems: testutils.Problems{},
		},
		{
			testName: "Invalid-MissingRequestsField",
			Field:    "string parent",
			problems: testutils.Problems{{Message: `no "requests" field`}},
		},
		{
			testName: "Invalid-RequestsFieldIsNotRepeated",
			Field:    "UpdateBookRequest requests",
			problems: testutils.Problems{{Message: "repeated"}},
		},
		{
			testName: "Invalid-RequestsFieldWrongType",
			Field:    "repeated int32 requests",
			problems: testutils.Problems{{
				Suggestion: "UpdateBookRequest",
			}},
		},
		{
			testName: "Invalid-RequestsNotRepeatedWrongType",
			Field:    "int32 requests",
			problems: testutils.Problems{
				{Message: `The "requests" field should be repeated`},
				{
					Suggestion: "UpdateBookRequest",
				}},
		},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
				message BatchUpdateBooksRequest {
					{{.Field}} = 1;
				}
				message UpdateBookRequest {}
				`, test)

			// Determine the descriptor that a failing test will attach to.
			var problemDesc desc.Descriptor
			if requests := file.GetMessageTypes()[0].FindFieldByName("requests"); requests != nil {
				problemDesc = requests
			} else {
				problemDesc = file.GetMessageTypes()[0]
			}

			problems := requestRequestsField.Lint(file)
			if diff := test.problems.SetDescriptor(problemDesc).Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
