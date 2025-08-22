// Copyright 2021 Google LLC
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

package aip0132

import (
	"testing"

	"github.com/googleapis/api-linter/v2/rules/internal/testutils"
)

func TestHTTPURIParent(t *testing.T) {
	tests := []struct {
		testName     string
		URI          string
		MethodName   string
		RequestField string
		problems     testutils.Problems
	}{
		{"Valid", "/v1/{parent=publishers/*/books/*}", "ListBooks", "string parent = 1;", nil},
		{"InvalidVarParent", "/v1/{book=publishers/*/books/*}", "ListBooks", "string parent = 1;", testutils.Problems{{Message: "HTTP URI should include a `parent` variable."}}},
		{"InvalidNoVarParent", "/v1/publishers/*/books/*", "ListBooks", "string parent = 1;", testutils.Problems{{Message: "HTTP URI should include a `parent` variable."}}},
		{"ValidNoParent", "/v1/books/*", "ListBooks", "", nil},
		{"Irrelevant", "/v1/{book=publishers/*/books/*}", "BuildBook", "string parent = 1;", nil},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/api/annotations.proto";
				service Library {
					rpc {{.MethodName}}({{.MethodName}}Request) returns ({{.MethodName}}Response) {
						option (google.api.http) = {
							get: "{{.URI}}"
						};
					}
				}
				message {{.MethodName}}Request {
					{{.RequestField}}
				}
				message {{.MethodName}}Response {}
			`, test)
			method := f.Services().Get(0).Methods().Get(0)
			if diff := test.problems.SetDescriptor(method).Diff(httpURIParent.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
