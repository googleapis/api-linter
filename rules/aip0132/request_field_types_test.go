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

package aip0132

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestRequestFieldTypes(t *testing.T) {
	// Set up the testing permutations.
	tests := []struct {
		testName string
		Message  string
		Field    string
		problems testutils.Problems
	}{
		{"Filter", "ListBooksRequest", "string filter", nil},
		{"FilterInvalid", "ListBooksRequest", "bytes filter", testutils.Problems{{Message: "singular string", Suggestion: "string"}}},
		{"FilterInvalidRepeated", "ListBooksRequest", "repeated string filter", testutils.Problems{{Message: "singular string", Suggestion: "string"}}},
		{"OrderBy", "ListBooksRequest", "string order_by", nil},
		{"OrderByInvalid", "ListBooksRequest", "bytes order_by", testutils.Problems{{Message: "singular string", Suggestion: "string"}}},
		{"OrderByInvalidRepeated", "ListBooksRequest", "repeated string order_by", testutils.Problems{{Message: "singular string", Suggestion: "string"}}},
		{"ShowDeleted", "ListBooksRequest", "bool show_deleted", nil},
		{"ShowDeletedInvalid", "ListBooksRequest", "int32 show_deleted", testutils.Problems{{Message: "singular bool", Suggestion: "bool"}}},
		{"ShowDeletedInvalidRepeated", "ListBooksRequest", "repeated bool show_deleted", testutils.Problems{{Message: "singular bool", Suggestion: "bool"}}},
		{"IrrelevantMessage", "Book", "bytes order_by", nil},
		{"IrrelevantField", "ListBooksRequest", "bytes foo", nil},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				message {{.Message}} {
					{{.Field}} = 1;
				}
			`, test)
			field := f.GetMessageTypes()[0].GetFields()[0]
			problems := requestFieldTypes.Lint(f)
			if diff := test.problems.SetDescriptor(field).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}
