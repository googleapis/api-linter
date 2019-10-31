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

package aip0122

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestHttpUriField(t *testing.T) {
	tests := []struct {
		testName string
		URI      string
		problems testutils.Problems
	}{
		{"Valid", "/v1/{name=publishers/*/books/*}:frob", testutils.Problems{}},
		{"ValidCamelPattern", "/v1/{name=publishers/*/frobbableBooks/*}:frob", testutils.Problems{}},
		{"InvalidSnakePattern", "/v1/{name=publishers/*/frobbable_books/*}:frob", testutils.Problems{{Message: "URI patterns"}}},
		{"InvalidCamelVariable", "/v1/{bookName=publishers/*/books/*}:frob", testutils.Problems{{Message: "Variable names"}}},
		{"ValidSnakeVariable", "/v1/{book_name=publishers/*/books/*}:frob", testutils.Problems{}},
		{"ValidSnakeSoloVariable", "/v1/{book_name}:frob", testutils.Problems{}},
		{"InvalidCamelSoloVariable", "/v1/{bookName}:frob", testutils.Problems{{Message: "Variable names"}}},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
				import "google/api/annotations.proto";
				service Library {
					rpc FrobBook(FrobBookRequest) returns (FrobBookResponse) {
						option (google.api.http) = {
							post: "{{.URI}}"
							body: "*"
						};
					}
				}
				message FrobBookRequest {}
				message FrobBookResponse {}
			`, test)
			method := file.GetServices()[0].GetMethods()[0]
			problems := httpURICase.Lint(file)
			if diff := test.problems.SetDescriptor(method).Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
