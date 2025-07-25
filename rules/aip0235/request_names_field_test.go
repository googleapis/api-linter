// Copyright 2020 Google LLC
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

package aip0235

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestNamesField(t *testing.T) {
	// Set up the testing permutations.
	tests := []struct {
		testName    string
		src         string
		problems    testutils.Problems
		problemDesc func(m protoreflect.MessageDescriptor) protoreflect.Descriptor
	}{
		{
			testName: "Valid-Names",
			src: `
				message BatchDeleteBooksRequest {
					repeated string names = 1;
				}
			`,
			problems: testutils.Problems{},
		},
		{
			testName: "Valid-StandardDeleteReq",
			src: `
				message BatchDeleteBooksRequest {
					repeated DeleteBookRequest requests = 1;
				}

				message DeleteBookRequest {}
			`,
			problems: testutils.Problems{},
		},
		{
			testName: "Invalid-MissingNamesField",
			src: `
				message BatchDeleteBooksRequest {
					string parent = 1;
				}
			`,
			problems: testutils.Problems{{Message: `Message "BatchDeleteBooksRequest" has no "names" field`}},
		},
		{
			testName: "Invalid-KeepingNamesFieldOnly",
			src: `
				message BatchDeleteBooksRequest {
					repeated string names = 1;
					repeated DeleteBookRequest requests = 2;
				}

				message DeleteBookRequest {}
			`,
			problems: testutils.Problems{{Message: `Message "BatchDeleteBooksRequest" should delete "requests" field, only keep the "names" field`}},
			problemDesc: func(m protoreflect.MessageDescriptor) protoreflect.Descriptor {
				return m.Fields().ByName("requests")
			},
		},
		{
			testName: "Invalid-NamesFieldIsNotRepeated",
			src: `
				message BatchDeleteBooksRequest {
					string names = 1;
				}`,
			problems: testutils.Problems{{Suggestion: "repeated string"}},
			problemDesc: func(m protoreflect.MessageDescriptor) protoreflect.Descriptor {
				return m.Fields().ByName("names")
			},
		},
		{
			testName: "Invalid-NamesFieldWrongType",
			src: `
				message BatchDeleteBooksRequest {
					repeated int32 names = 1;
				}
			`,
			problems: testutils.Problems{{Suggestion: "string"}},
			problemDesc: func(m protoreflect.MessageDescriptor) protoreflect.Descriptor {
				return m.Fields().ByName("names")
			},
		},
		{
			testName: "Invalid-DeleteReqFieldIsNotRepeated",
			src: `
				message BatchDeleteBooksRequest {
					DeleteBookRequest requests = 1;
				}

				message DeleteBookRequest {}
			`,
			problems: testutils.Problems{{Message: `The "requests" field should be repeated`}},
			problemDesc: func(m protoreflect.MessageDescriptor) protoreflect.Descriptor {
				return m.Fields().ByName("requests")
			},
		},
		{
			testName: "Invalid-DeleteReqFieldWrongType",
			src: `
				message BatchDeleteBooksRequest {
					repeated string requests = 1;
				}
			`,
			problems: testutils.Problems{{Message: `The "requests" field on Batch Delete Request should be a "DeleteBookRequest" type`}},
			problemDesc: func(m protoreflect.MessageDescriptor) protoreflect.Descriptor {
				return m.Fields().ByName("requests")
			},
		},
		{
			testName: "Irrelevant-UnmatchedMessageName",
			src:      `message DeleteBooksRequest {}`,
			problems: testutils.Problems{},
		},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			file := testutils.ParseProto3String(t, test.src)

			m := file.Messages().Get(0)

			// Determine the descriptor that a failing test will attach to.
			var problemDesc protoreflect.Descriptor = m
			if test.problemDesc != nil {
				problemDesc = test.problemDesc(m)
			}

			problems := requestNamesField.Lint(file)
			if diff := test.problems.SetDescriptor(problemDesc).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}
