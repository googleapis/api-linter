// Copyright 2020 Google LLC
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

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestRequestMessageName(t *testing.T) {
	// Set up the testing permutations.
	tests := []struct {
		testName       string
		MethodName     string
		ReqMessageName string
		problems       testutils.Problems
	}{
		{"Valid", "PurgeBooks", "PurgeBooksRequest", nil},
		{"Invalid", "PurgeBooks", "PurgeBookRequest", testutils.Problems{{Suggestion: "PurgeBooksRequest"}}},
		{"Irrelevant", "AcquireBook", "Book", nil},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				service Library {
					rpc {{.MethodName}}({{.ReqMessageName}}) returns ({{.MethodName}}Response) {}
				}
				message {{.ReqMessageName}} {}
				message {{.MethodName}}Response {}
			`, test)
			m := f.Services()[0].Methods()[0]
			if diff := test.problems.SetDescriptor(m).Diff(requestMessageName.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
