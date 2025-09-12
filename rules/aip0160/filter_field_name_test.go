// Copyright 2025 Google LLC
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

package aip0160

import (
	"fmt"
	"testing"

	"github.com/googleapis/api-linter/v2/lint"
	"github.com/googleapis/api-linter/v2/rules/internal/testutils"
)

func TestFiltersFieldName(t *testing.T) {
	tests := []struct {
		testName   string
		FieldType  string
		FieldName  string
		MethodName string
		wantErr    bool
	}{
		{"Valid", "string", "filter", "ListBooks", false},
		{"Valid different native type List method", "bytes", "filter", "ListBooks", false},
		{"Valid different message type List method", "BookFilters", "filter", "ListBooks", false},
		{"Valid different native type custom method", "bytes", "filter", "SearchBooks", false},
		{"Valid different message type custom method", "BookFilters", "filter", "SearchBooks", false},
		{"Invalid list method", "string", "filters", "ListBooks", true},
		{"Invalid custom method", "string", "filters", "SearchBooks", true},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
				service Library {
					rpc {{.MethodName}}({{.MethodName}}Request) returns ({{.MethodName}}Response) {}
					rpc GetBook(GetBookRequest) returns (Book) {}
				}
				message {{.MethodName}}Request {
					{{.FieldType}} {{.FieldName}} = 1;
				}
				message {{.MethodName}}Response {}

				message GetBookRequest {}
				message Book { string filters = 1; }

				message BookFilters {}
		`, test)
			field := file.Messages().Get(0).Fields().Get(0)
			wantProblems := testutils.Problems{}
			if test.wantErr {
				wantProblems = append(wantProblems, lint.Problem{
					Message:    fmt.Sprintf(`"string filter", not "%s %s`, test.FieldType, test.FieldName),
					Suggestion: "filter",
					Descriptor: field,
				})
			}
			if diff := wantProblems.Diff(filterFieldName.Lint(file)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
