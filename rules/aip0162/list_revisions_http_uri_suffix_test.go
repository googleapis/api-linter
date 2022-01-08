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

package aip0162

import (
	"testing"

	"github.com/commure/api-linter/rules/internal/testutils"
)

func TestListRevisionsHTTPURISuffix(t *testing.T) {
	tests := []struct {
		testName   string
		URI        string
		MethodName string
		problems   testutils.Problems
	}{
		{"Valid", "/v1/{name=publishers/*/books/*}:listRevisions", "ListBookRevisions", nil},
		{"InvalidSuffix", "/v1/{name=publishers/*/books/*}:list", "ListBookRevisions", testutils.Problems{{Message: ":listRevisions"}}},
		{"Irrelevant", "/v1/{name=publishers/*/books/*}:list", "ListBooks", nil},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
				import "google/api/annotations.proto";
				service Library {
					rpc {{.MethodName}}({{.MethodName}}Request) returns ({{.MethodName}}Response) {
						option (google.api.http) = {
							get: "{{.URI}}"
						};
					}
				}
				message {{.MethodName}}Request {}
				message {{.MethodName}}Response {}
				`, test)

			problems := listRevisionsHTTPURISuffix.Lint(file)
			if diff := test.problems.SetDescriptor(file.GetServices()[0].GetMethods()[0]).Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
