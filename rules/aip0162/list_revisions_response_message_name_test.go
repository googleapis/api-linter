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

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestListRevisionsResponseMessageName(t *testing.T) {
	for _, test := range []struct {
		name         string
		Method       string
		ResponseType string
		problems     testutils.Problems
	}{
		{"Valid", "ListBookRevisions", "ListBookRevisionsResponse", nil},
		{"Invalid", "ListBookRevisions", "ListRevisionsResponse", testutils.Problems{{Suggestion: "ListBookRevisionsResponse"}}},
		{"Irrelevant", "ListBooks", "Book", nil},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				service Library {
					rpc {{.Method}}({{.Method}}Request) returns ({{.ResponseType}});
				}
				message {{.Method}}Request {}
				message {{.ResponseType}} {}
			`, test)
			m := f.GetServices()[0].GetMethods()[0]
			if diff := test.problems.SetDescriptor(m).Diff(listRevisionsResponseMessageName.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
