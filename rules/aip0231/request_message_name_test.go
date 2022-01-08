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

	"github.com/commure/api-linter/rules/internal/testutils"
)

func TestInputName(t *testing.T) {
	// Set up the testing permutations.
	tests := []struct {
		testName   string
		MethodName string
		InputName  string
		problems   testutils.Problems
	}{
		{"Valid", "BatchGetBooks", "BatchGetBooksRequest", testutils.Problems{}},
		{"Invalid", "BatchGetBooks", "GetBookRequest", testutils.Problems{{Suggestion: "BatchGetBooksRequest"}}},
		{"Irrelevant", "GetBookRequest", "Foo", testutils.Problems{}},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
				import "google/api/annotations.proto";

				service Library {
					rpc {{.MethodName}}({{.InputName}}) returns ({{.MethodName}}Response) {
						option (google.api.http) = {
							get: "/v1/{parent=publishers/*}/books:batchGet"
						};
					}
				}

				message {{.InputName}} {}
				message {{.MethodName}}Response {}
			`, test)

			m := file.GetServices()[0].GetMethods()[0]

			problems := inputName.Lint(file)
			if diff := test.problems.SetDescriptor(m).Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
