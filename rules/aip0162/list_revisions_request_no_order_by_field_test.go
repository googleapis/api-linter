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

func TestListRevisionsRequestNoOrderByField(t *testing.T) {
	for _, test := range []struct {
		name     string
		Message  string
		Field    string
		problems testutils.Problems
	}{
		{"Valid", "ListBookRevisionsRequest", "string name = 1;", nil},
		{"Invalid", "ListBookRevisionsRequest", "string order_by = 1;", testutils.Problems{{Message: "not contain an `order_by`"}}},
		{"IrrelevantMessage", "ListBooksRequest", "string order_by = 1;", nil},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				message {{.Message}} {
					{{.Field}}
				}
			`, test)
			d := f.GetMessageTypes()[0].GetFields()[0]
			if diff := test.problems.SetDescriptor(d).Diff(listRevisionsRequestNoOrderByField.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
