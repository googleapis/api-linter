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

package aip0233

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
	"github.com/jhump/protoreflect/desc"
)

func TestRequestRequestsField(t *testing.T) {
	// Set up the testing permutations.
	tests := []struct {
		testName    string
		field       string
		problems    testutils.Problems
		problemDesc func(m *desc.MessageDescriptor) desc.Descriptor
	}{
		{
			testName: "Valid",
			field:    "repeated CreateBookRequest requests = 1;",
			problems: testutils.Problems{},
		},
		{
			testName: "Invalid-MissingRequestsField",
			field:    "string parent = 1;",
			problems: testutils.Problems{{Message: `no "requests" field`}},
		},
		{
			testName: "Invalid-RequestsFieldIsNotRepeated",
			field:    "CreateBookRequest requests = 1;",
			problems: testutils.Problems{{Message: `The "requests" field should be repeated`}},
			problemDesc: func(m *desc.MessageDescriptor) desc.Descriptor {
				return m.FindFieldByName("requests")
			},
		},
		{
			testName: "Invalid-RequestsFieldWrongType",
			field:    "repeated int32 requests = 1;",
			problems: testutils.Problems{{
				Message:    `should be a "CreateBookRequest" type`,
				Suggestion: "CreateBookRequest",
			}},
			problemDesc: func(m *desc.MessageDescriptor) desc.Descriptor {
				return m.FindFieldByName("requests")
			},
		},
		{
			testName: "Invalid-RequestsNotRepeatedWrongType",
			field:    "int32 requests = 1;",
			problems: testutils.Problems{
				{Message: `The "requests" field should be repeated`},
				{
					Message:    `should be a "CreateBookRequest" type`,
					Suggestion: "CreateBookRequest",
				}},
			problemDesc: func(m *desc.MessageDescriptor) desc.Descriptor {
				return m.FindFieldByName("requests")
			},
		},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			template := `
message BatchCreateBooksRequest {
	{{.Field}}
}

message CreateBookRequest {}
`
			file := testutils.ParseProto3Tmpl(t, template,
				struct{ Field string }{test.field})

			m := file.GetMessageTypes()[0]

			// Determine the descriptor that a failing test will attach to.
			var problemDesc desc.Descriptor = m
			if test.problemDesc != nil {
				problemDesc = test.problemDesc(m)
			}

			problems := requestRequestsField.Lint(file)
			if diff := test.problems.SetDescriptor(problemDesc).Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
