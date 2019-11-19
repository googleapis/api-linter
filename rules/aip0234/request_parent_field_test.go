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

func TestRequestParentField(t *testing.T) {
	// Set up the testing permutations.
	tests := []struct {
		testName    string
		MessageName string
		Field       string
		problems    testutils.Problems
	}{
		{
			"Valid",
			"BatchUpdateBooksRequest",
			"string parent",
			testutils.Problems{},
		},
		{
			"MissingField",
			"BatchUpdateBooksRequest",
			"string id",
			testutils.Problems{{Message: "parent"}},
		},
		{
			"InvalidType",
			"BatchUpdateBooksRequest",
			"int32 parent",
			testutils.Problems{{
				Message:    "string",
				Suggestion: "string",
			}},
		},
		{
			"Irrelevant",
			"EnumerateBooksRequest",
			"string id",
			testutils.Problems{},
		},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
				message {{.MessageName}} {
					{{.Field}} = 1;
					repeated UpdateBookRequest requests = 2;
				}
				message UpdateBookRequest{}
				`, test)

			// Determine the descriptor that a failing test will attach to.
			var problemDesc desc.Descriptor
			if parent := file.GetMessageTypes()[0].FindFieldByName("parent"); parent != nil {
				problemDesc = parent
			} else {
				problemDesc = file.GetMessageTypes()[0]
			}

			// Run the lint rule, and establish that it returns the correct problems.
			problems := requestParentField.Lint(file)
			if diff := test.problems.SetDescriptor(problemDesc).Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
