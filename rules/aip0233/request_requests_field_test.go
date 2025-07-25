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
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestRequestRequestsField(t *testing.T) {
	// Set up the testing permutations.
	tests := []struct {
		testName    string
		Field       string
		problems    testutils.Problems
		problemDesc func(m protoreflect.MessageDescriptor) protoreflect.Descriptor
	}{
		{
			testName: "Valid",
			Field:    "repeated CreateBookRequest requests = 1;",
			problems: testutils.Problems{},
		},
		{
			testName: "Invalid-MissingRequestsField",
			Field:    "string parent = 1;",
			problems: testutils.Problems{{Message: `no "requests" field`}},
		},
		{
			testName: "Invalid-RequestsFieldIsNotRepeated",
			Field:    "CreateBookRequest requests = 1;",
			problems: testutils.Problems{{Message: "repeated"}},
			problemDesc: func(m protoreflect.MessageDescriptor) protoreflect.Descriptor {
				return m.Fields().ByName("requests")
			},
		},
		{
			testName: "Invalid-RequestsFieldWrongType",
			Field:    "repeated int32 requests = 1;",
			problems: testutils.Problems{{
				Suggestion: "CreateBookRequest",
			}},
			problemDesc: func(m protoreflect.MessageDescriptor) protoreflect.Descriptor {
				return m.Fields().ByName("requests")
			},
		},
		{
			testName: "Invalid-RequestsNotRepeatedWrongType",
			Field:    "int32 requests = 1;",
			problems: testutils.Problems{
				{Message: "repeated"},
				{
					Suggestion: "CreateBookRequest",
				},
			},
			problemDesc: func(m protoreflect.MessageDescriptor) protoreflect.Descriptor {
				return m.Fields().ByName("requests")
			},
		},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
				message BatchCreateBooksRequest {
					{{.Field}}
				}
				
				message CreateBookRequest {}
				`, test)

			m := file.Messages().Get(0)

			// Determine the descriptor that a failing test will attach to.
			var problemDesc protoreflect.Descriptor = m
			if test.problemDesc != nil {
				problemDesc = test.problemDesc(m)
			}

			problems := requestRequestsField.Lint(file)
			if diff := test.problems.SetDescriptor(problemDesc).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}
