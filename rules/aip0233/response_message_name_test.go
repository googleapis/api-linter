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
)

func TestResponseMessageName(t *testing.T) {
	// Set up the testing permutations.
	tests := []struct {
		testName   string
		MethodName string
		Response   string
		problems   testutils.Problems
	}{
		{
			testName:   "Valid-BatchCreateBooksResponse",
			MethodName: "BatchCreateBooks",
			Response:   "BatchCreateBooksResponse",
			problems:   testutils.Problems{},
		},
		{
			testName:   "Invalid-MissMatchingMethodName",
			MethodName: "BatchCreateBooks",
			Response:   "BatchCreateBookResponse",
			problems: testutils.Problems{{
				Suggestion: "BatchCreateBooksResponse",
			}},
		},
		{
			testName:   "Irrelevant",
			MethodName: "CreateBook",
			Response:   "Book",
			problems:   testutils.Problems{},
		},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
				import "google/api/annotations.proto";
				
				service BookService {
					rpc {{.MethodName}}({{.MethodName}}Request) returns ({{.Response}}) {
						option (google.api.http) = {
							post: "/v1/{parent=publishers/*}/"
							body: "*"
						};
					}
				}
				
				message {{.MethodName}}Request {}
				
				message {{.Response}}{}
				`, test)

			m := file.GetServices()[0].GetMethods()[0]

			problems := responseMessageName.Lint(file)
			if diff := test.problems.SetDescriptor(m).Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

// Other cases for long running response will be handled by AIP-0151
func TestLongRunningResponse(t *testing.T) {
	// Set up the testing permutations.
	tests := []struct {
		testName     string
		ResponseType string
		problems     testutils.Problems
	}{
		{
			testName:     "Valid-LongRunning",
			ResponseType: "BatchCreateBooksResponse",
			problems:     testutils.Problems{},
		},
		{
			testName: "Valid-LongRunningEmptyResponseType",
			problems: testutils.Problems{},
		},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
				import "google/api/annotations.proto";
				import "google/longrunning/operations.proto";
				
				service BookService {
					rpc BatchCreateBooks(BatchCreateBooksRequest) returns (google.longrunning.Operation) {
						option (google.api.http) = {
							post: "/v1/{parent=publishers/*}/books:batchCreate"
							body: "*"
						};
						option (google.longrunning.operation_info) = {
							response_type: "{{.ResponseType}}"
						};
					}
				}
				
				message BatchCreateBooksRequest {}
				
				message BatchCreateBooksResponse{}
				`, test)

			m := file.GetServices()[0].GetMethods()[0]

			problems := responseMessageName.Lint(file)
			if diff := test.problems.SetDescriptor(m).Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
