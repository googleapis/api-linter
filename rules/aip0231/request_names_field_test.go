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

package aip0231

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
	"github.com/jhump/protoreflect/desc"
)

func TestNamesField(t *testing.T) {
	// Set up the testing permutations.
	tests := []struct {
		testName    string
		src         string
		problems    testutils.Problems
		problemDesc func(m *desc.MessageDescriptor) desc.Descriptor
	}{
		{
			testName: "Valid-Names",
			src: `
				message BatchGetBooksRequest {
					repeated string names = 1;
				}
			`,
			problems: testutils.Problems{},
		},
		{
			testName: "Valid-StandardGetReq",
			src: `
				message BatchGetBooksRequest {
				repeated GetBookRequest requests = 1;
				}

				message GetBookRequest {}
			`,
			problems: testutils.Problems{},
		},
		{
			testName: "Invalid-MissingNamesField",
			src: `
				message BatchGetBooksRequest {
				string parent = 1;
				}
			`,
			problems: testutils.Problems{{Message: `Message "BatchGetBooksRequest" has no "names" field`}},
		},
		{
			testName: "Invalid-KeepingNamesFieldOnly",
			src: `
				message BatchGetBooksRequest {
				repeated string names = 1;
				repeated GetBookRequest requests = 2;
				}

				message GetBookRequest {}
			`,
			problems: testutils.Problems{{Message: `Message "BatchGetBooksRequest" should delete "requests" field, only keep the "names" field`}},
			problemDesc: func(m *desc.MessageDescriptor) desc.Descriptor {
				return m.FindFieldByName("requests")
			},
		},
		{
			testName: "Invalid-NamesFieldIsNotRepeated",
			src: `
				message BatchGetBooksRequest {
					string names = 1;
				}`,
			problems: testutils.Problems{{Suggestion: "repeated string"}},
			problemDesc: func(m *desc.MessageDescriptor) desc.Descriptor {
				return m.FindFieldByName("names")
			},
		},
		{
			testName: "Invalid-NamesFieldWrongType",
			src: `
				message BatchGetBooksRequest {
					repeated int32 names = 1;
				}
			`,
			problems: testutils.Problems{{Suggestion: "string"}},
			problemDesc: func(m *desc.MessageDescriptor) desc.Descriptor {
				return m.FindFieldByName("names")
			},
		},
		{
			testName: "Invalid-GetReqFieldIsNotRepeated",
			src: `
				message BatchGetBooksRequest {
				GetBookRequest requests = 1;
				}

				message GetBookRequest {}
			`,
			problems: testutils.Problems{{Message: `The "requests" field should be repeated`}},
			problemDesc: func(m *desc.MessageDescriptor) desc.Descriptor {
				return m.FindFieldByName("requests")
			},
		},
		{
			testName: "Invalid-NamesFieldWrongType",
			src: `
				message BatchGetBooksRequest {
					repeated string requests = 1;
				}
			`,
			problems: testutils.Problems{{Message: `The "requests" field on Batch Get Request should be a "GetBookRequest" type`}},
			problemDesc: func(m *desc.MessageDescriptor) desc.Descriptor {
				return m.FindFieldByName("requests")
			},
		},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			file := testutils.ParseProto3String(t, test.src)

			m := file.GetMessageTypes()[0]

			// Determine the descriptor that a failing test will attach to.
			var problemDesc desc.Descriptor = m
			if test.problemDesc != nil {
				problemDesc = test.problemDesc(m)
			}

			problems := namesField.Lint(file)
			if diff := test.problems.SetDescriptor(problemDesc).Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
