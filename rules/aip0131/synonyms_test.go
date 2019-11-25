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

package aip0131

import (
	"testing"

	"github.com/googleapis/api-linter/descrule"
	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestSynonyms(t *testing.T) {
	tests := []struct {
		MethodName string
		problems   testutils.Problems
	}{
		{"GetBook", testutils.Problems{}},
		{"AcquireBook", testutils.Problems{{Suggestion: "GetBook"}}},
		{"FetchBook", testutils.Problems{{Suggestion: "GetBook"}}},
		{"LookupBook", testutils.Problems{{Suggestion: "GetBook"}}},
		{"ReadBook", testutils.Problems{{Suggestion: "GetBook"}}},
		{"RetrieveBook", testutils.Problems{{Suggestion: "GetBook"}}},
	}
	for _, test := range tests {
		file := testutils.ParseProto3Tmpl(t, `
			service Library {
				rpc {{.MethodName}}({{.MethodName}}Request) returns (Book);
			}
			message {{.MethodName}}Request {}
			message Book {}
		`, test)
		m := file.GetServices()[0].GetMethods()[0]
		if diff := test.problems.SetDescriptor(m).Diff(synonyms.Lint(descrule.NewMethod(m))); diff != "" {
			t.Errorf(diff)
		}
	}
}
