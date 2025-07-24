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
	"strings"
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

type responseMessageNameTest struct {
	testName         string
	MethodName       string
	MethodAnnotation string
	Response         string
	problems         testutils.Problems
}

func (t responseMessageNameTest) ResponseExternal() bool {
	return strings.Contains(t.Response, ".")
}

func TestResponseMessageName(t *testing.T) {
	// Set up the testing permutations.
	tests := []responseMessageNameTest{
		{
			testName:   "Valid-BatchDeleteBooksResponse",
			MethodName: "BatchDeleteBooks",
			Response:   "BatchDeleteBooksResponse",
			problems:   testutils.Problems{},
		},
		{
			testName:   "Valid-Empty",
			MethodName: "BatchDeleteBooks",
			Response:   "google.protobuf.Empty",
			problems:   testutils.Problems{},
		},
		{
			testName:   "Valid-LRO-BatchDeleteBooksResponse",
			MethodName: "BatchDeleteBooks",
			Response:   "google.longrunning.Operation",
			MethodAnnotation: `option (google.longrunning.operation_info) = {
				response_type: "BatchDeleteBooksResponse"
			};`,
			problems: testutils.Problems{},
		},
		{
			testName:   "Valid-LRO-Empty",
			MethodName: "BatchDeleteBooks",
			Response:   "google.longrunning.Operation",
			MethodAnnotation: `option (google.longrunning.operation_info) = {
				response_type: "google.protobuf.Empty"
			};`,
			problems: testutils.Problems{},
		},
		{
			testName:   "Invalid",
			MethodName: "BatchDeleteBooks",
			Response:   "BatchDeleteBookResponse",
			problems: testutils.Problems{{
				Message:    "`BatchDeleteBooksResponse`",
				Suggestion: "google.protobuf.Empty",
			}},
		},
		{
			testName:   "Invalid-LRO",
			MethodName: "BatchDeleteBooks",
			Response:   "google.longrunning.Operation",
			MethodAnnotation: `option (google.longrunning.operation_info) = {
				response_type: "BatchDeleteBookResponse"
			};`,
			problems: testutils.Problems{{
				Message:    "`BatchDeleteBooksResponse`",
				Suggestion: "google.protobuf.Empty",
			}},
		},
		{
			testName:   "Irrelevant",
			MethodName: "DeleteBook",
			Response:   "Book",
			problems:   testutils.Problems{},
		},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
				import "google/protobuf/empty.proto";
				import "google/longrunning/operations.proto";

				service BookService {
					rpc {{.MethodName}}({{.MethodName}}Request) returns ({{.Response}}) {
						{{.MethodAnnotation}}
					}
				}
				message {{.MethodName}}Request {}
				{{ if not .ResponseExternal }}message {{.Response}} {}{{ end }}
				`, test)

			m := file.Services()[0].Methods()[0]

			problems := responseMessageName.Lint(file)
			if diff := test.problems.SetDescriptor(m).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}
