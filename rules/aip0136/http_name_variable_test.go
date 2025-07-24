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

func TestHTTPVariables(t *testing.T) {
	for _, test := range []struct {
		name       string
		MethodName string
		URI        string
		problems   testutils.Problems
	}{
		{"Valid", "WriteBook", "/v1/{name=publishers/*/books/*}:write", testutils.Problems{}},
		{"ValidRevision", "TagBookRevision", "/v1/{name=publishers/*/books/*}:tagRevision", testutils.Problems{}},
		{"ValidRevisions", "ListBookRevisions", "/v1/{name=publishers/*/books/*}:listRevisions", testutils.Problems{}},
		{"Invalid", "WritePage", "/v1/{name=publishers/*/books/*}:writePage", testutils.Problems{{Message: "name variable"}}},
		{"ValidBookVar", "WritePage", "/v1/{book=publishers/*/books/*}:writePage", testutils.Problems{}},
		{"ValidTwoWordNoun", "WriteAudioBook", "/v1/{name=publishers/*/audioBooks/*}:write", testutils.Problems{}},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/api/annotations.proto";
				service Library {
					rpc {{.MethodName}}({{.MethodName}}Request)
							returns ({{.MethodName}}Response) {
						option (google.api.http) = {
							post: "{{.URI}}"
							body: "*"
						};
					}
				}
				message {{.MethodName}}Request {}
				message {{.MethodName}}Response {}
			`, test)
			m := f.Services()[0].Methods()[0]
			if diff := test.problems.SetDescriptor(m).Diff(httpNameVariable.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
