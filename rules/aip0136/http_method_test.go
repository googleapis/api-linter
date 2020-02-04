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

package aip0136

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestHttpMethod(t *testing.T) {
	tests := []struct {
		testName   string
		MethodName string
		HTTPMethod string
		problems   testutils.Problems
	}{
		{"ValidGet", "ArchiveBook", "get", testutils.Problems{}},
		{"ValidPost", "ArchiveBook", "post", testutils.Problems{}},
		{"InvalidPut", "ArchiveBook", "put", testutils.Problems{{Message: "POST or GET"}}},
		{"InvalidPatch", "ArchiveBook", "patch", testutils.Problems{{Message: "POST or GET"}}},
		{"InvalidDelete", "ArchiveBook", "delete", testutils.Problems{{Message: "POST or GET"}}},
		{"IrrelevantPatch", "UpdateBook", "patch", testutils.Problems{}},
		{"IrrelevantDelete", "DeleteBook", "delete", testutils.Problems{}},
		{"IrrelevantPut", "ReplaceBook", "put", testutils.Problems{}},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
				import "google/api/annotations.proto";
				service Library {
					rpc {{.MethodName}}({{.MethodName}}Request) returns ({{.MethodName}}Response) {
						option (google.api.http) = {
							{{.HTTPMethod}}: "/v1/{{.MethodName}}"
						};
					}
				}
				message {{.MethodName}}Request {}
				message {{.MethodName}}Response {}
			`, test)
			method := file.GetServices()[0].GetMethods()[0]
			got := httpMethod.Lint(file)
			if diff := test.problems.SetDescriptor(method).Diff(got); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
