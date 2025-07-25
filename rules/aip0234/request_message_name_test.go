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
)

func TestRequestMessageName(t *testing.T) {
	// Set up the testing permutations.
	tests := []struct {
		testName   string
		MethodName string
		Request    string
		problems   testutils.Problems
	}{
		{
			testName:   "Valid-BatchUpdateBooksRequest",
			MethodName: "BatchUpdateBooks",
			Request:    "BatchUpdateBooksRequest",
			problems:   testutils.Problems{},
		},
		{
			testName:   "Invalid-MissMatchingMethodName",
			MethodName: "BatchUpdateBooks",
			Request:    "BatchUpdateBookRequest",
			problems: testutils.Problems{{
				Suggestion: "BatchUpdateBooksRequest",
			}},
		},
		{
			testName:   "Irrelevant",
			MethodName: "AcquireBook",
			Request:    "AcquireBookRequest",
			problems:   testutils.Problems{},
		},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
				import "google/api/annotations.proto";
				
				service BookService {
					rpc {{.MethodName}}({{.Request}}) returns ({{.MethodName}}Response) {
						option (google.api.http) = {
							post: "/v1/{parent=publishers/*}/"
							body: "*"
						};
					}
				}
				message {{.Request}}{}
				message {{.MethodName}}Response{}
				`, test)

			m := file.Services().Get(0).Methods().Get(0)

			problems := requestMessageName.Lint(file)
			if diff := test.problems.SetDescriptor(m).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}
