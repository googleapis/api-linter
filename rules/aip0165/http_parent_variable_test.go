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

package aip0165

import (
	"testing"

	"github.com/googleapis/api-linter/v2/rules/internal/testutils"
)

func TestHTTPParentVariable(t *testing.T) {
	tests := []struct {
		testName   string
		URI        string
		MethodName string
		Field      string
		problems   testutils.Problems
	}{
		{"Valid", "/v1/{parent=publishers/*/books/*}", "PurgeBooks", "string parent = 1;", nil},
		{"InvalidVarParent", "/v1/{book=publishers/*/books/*}", "PurgeBooks", "string parent = 1;", testutils.Problems{{Message: "`parent`"}}},
		{"NoVarParent", "/v1/publishers/*/books/*", "PurgeBooks", "string parent = 1;", testutils.Problems{{Message: "`parent`"}}},
		{"NoParent", "/v1/books/*", "PurgeBooks", "", nil},
		{"Irrelevant", "/v1/{book=publishers/*/books/*}", "BuildBook", "string parent = 1;", nil},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/api/annotations.proto";
				  service Library {
					rpc {{.MethodName}}({{.MethodName}}Request) returns ({{.MethodName}}Response) {
						option (google.api.http) = {
							post: "{{.URI}}"
						};
					}
				}
				message {{.MethodName}}Request {
					{{.Field}}
				}
				message {{.MethodName}}Response {}
			`, test)
			method := f.Services().Get(0).Methods().Get(0)
			if diff := test.problems.SetDescriptor(method).Diff(httpParentVariable.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
