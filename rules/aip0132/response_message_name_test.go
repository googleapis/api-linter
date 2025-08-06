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

	"github.com/googleapis/api-linter/v2/rules/internal/testutils"
)

func TestResponseMessageName(t *testing.T) {
	// Set up the testing permutations.
	tests := []struct {
		testName        string
		MethodName      string
		RespMessageName string
		problems        testutils.Problems
	}{
		{"Valid", "ListBooks", "ListBooksResponse", testutils.Problems{}},
		{"Invalid", "ListBooks", "Books", testutils.Problems{{Suggestion: "ListBooksResponse"}}},
		{"Irrelevant", "EnumerateBooks", "EnumerateBooksResponse", testutils.Problems{}},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
				service Library {
					rpc {{.MethodName}}({{.MethodName}}Request) returns ({{.RespMessageName}}) {}
				}
				message {{.MethodName}}Request {}
				message {{.RespMessageName}} {}
			`, test)
			method := file.Services().Get(0).Methods().Get(0)
			problems := responseMessageName.Lint(file)
			if diff := test.problems.SetDescriptor(method).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}
