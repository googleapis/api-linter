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

package aip0133

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestHTTPURIParent(t *testing.T) {
	tests := []struct {
		testName   string
		URI        string
		MethodName string
		Pattern    string
		problems   testutils.Problems
	}{
		{"Valid", "/v1/{parent=publishers/*/books/*}", "CreateBook", "publishers/{publisher}/books/{book}", nil},
		{"InvalidVarParent", "/v1/{book=publishers/*/books/*}", "CreateBook", "publishers/{publisher}/books/{book}", testutils.Problems{{Message: "`parent` variable"}}},
		{"NoVarParent", "/v1/publishers/*/books/*", "CreateBook", "publishers/{publisher}/books/{book}", testutils.Problems{{Message: "`parent` variable"}}},
		{"NoParent", "/v1/books/*", "CreateBook", "books/{book}", nil},
		{"MultipleVars", "/v1/{parent=publishers/*}/{book=books/*}", "CreateBook", "publishers/{publisher}/books/{book}", testutils.Problems{{Message: "1 variable"}}},
		{"Irrelevant", "/v1/{book=publishers/*/books/*}", "BuildBook", "publishers/{publisher}/books/{book}", nil},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/api/annotations.proto";
				import "google/api/resource.proto";
				service Library {
					rpc {{.MethodName}}({{.MethodName}}Request) returns (Book) {
						option (google.api.http) = {
							post: "{{.URI}}"
						};
					}
				}
				message {{.MethodName}}Request {}
				message Book {
					option (google.api.resource) = {
						pattern: "{{.Pattern}}"
					};
				}
			`, test)
			method := f.GetServices()[0].GetMethods()[0]
			if diff := test.problems.SetDescriptor(method).Diff(httpURIParent.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
