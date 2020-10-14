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

package aip0133

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestInputName(t *testing.T) {
	// Set up the testing permutations.
	tests := []struct {
		testName   string
		MethodName string
		InputName  string
		problems   testutils.Problems
	}{
		{"Valid", "CreateBook", "CreateBookRequest", nil},
		{"Invalid", "CreateBook", "Book",
			testutils.Problems{{
				Message:    `Post RPCs should have a properly named request message "CreateBookRequest", but not "Book"`,
				Suggestion: "CreateBookRequest",
			}}},
		{"Irrelevant_OutputWrong", "CreateIamPolicy", "CreateIamPolicyRequest", nil},
		{"Irrelevant_NotCreate", "BuildBook", "Book", nil},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				service Library {
					rpc {{.MethodName}}({{.InputName}}) returns (Book) {}
				}
				message {{.InputName}} {}
				{{if ne .InputName "Book"}}
				message Book {}
				{{end}}
			`, test)
			m := f.GetServices()[0].GetMethods()[0]
			if diff := test.problems.SetDescriptor(m).Diff(inputName.Lint(f)); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
