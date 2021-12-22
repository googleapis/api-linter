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

package aip0134

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestHttpNameField(t *testing.T) {
	tests := []struct {
		testName   string
		URI        string
		MethodName string
		problems   testutils.Problems
	}{
		{"Valid", "/v1/{big_book.resource_name=publishers/*/books/*}", "UpdateBigBook", nil},
		{"ValidWithNumber", "/v1/{dv360.resource_name=publishers/*/dv360s/*}", "UpdateDv360", nil},
		{"InvalidNoUnderscore", "/v1/{bigbook.resource_name=publishers/*/books/*}",
			"UpdateBigBook", testutils.Problems{{Message: "`big_book.resource_name`"}}},
		{"InvalidVarNameBook", "/v1/{big_book=publishers/*/books/*}",
			"UpdateBigBook", testutils.Problems{{Message: "`big_book.resource_name`"}}},
		{"InvalidVarNameName", "/v1/{name=publishers/*/books/*}",
			"UpdateBigBook", testutils.Problems{{Message: "`big_book.resource_name`"}}},
		{"InvalidVarNameReversed", "/v1/{name.big_book=publishers/*/books/*}",
			"UpdateBigBook", testutils.Problems{{Message: "`big_book.resource_name`"}}},
		{"NoVarName", "/v1/publishers/*/books/*",
			"UpdateBigBook", testutils.Problems{{Message: "`big_book.resource_name`"}}},
		{"Irrelevant", "/v1/{book=publishers/*/books/*}",
			"AcquireBigBook", nil},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/api/annotations.proto";
				service Library {
					rpc {{.MethodName}}({{.MethodName}}Request) returns (BigBook) {
						option (google.api.http) = {
							patch: "{{.URI}}"
						};
					}
				}
				message BigBook {}
				message {{.MethodName}}Request {}
			`, test)
			method := f.GetServices()[0].GetMethods()[0]
			if diff := test.problems.SetDescriptor(method).Diff(httpNameField.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
