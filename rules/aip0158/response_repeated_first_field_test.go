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

package aip0158

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestResponseRepeatedFirstField(t *testing.T) {
	tests := []struct {
		testName    string
		FirstField  string
		SecondField string
		problems    testutils.Problems
	}{
		{"Valid", "repeated Book books = 1;", "string next_page_token = 2;", nil},
		{"SkippedZeroFields", "", "", nil},
		{"InvalidFirstFieldsNoMatch", "string next_page_token = 2;", "repeated Book books = 1;", testutils.Problems{{Message: "protobuf ID"}}},
		{"InvalidNotRepeated", "Book book = 1;", "string next_page_token = 2;", testutils.Problems{{Message: "repeated"}}},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Create the proto message.
			f := testutils.ParseProto3Tmpl(t, `
				message Book {
					string name = 1;
				}
				message ListBooksResponse {
					{{.FirstField}}
					{{.SecondField}}
				}
			`, test)

			// Determine the descriptor that a failing test will attach to.
			if m := f.GetMessageTypes()[1]; len(m.GetFields()) > 0 {
				test.problems.SetDescriptor(m.GetFields()[0])
			}

			// Run the lint rule and establish we get the correct problems.
			problems := responseRepeatedFirstField.Lint(f)
			if diff := test.problems.Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}

}
