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

package aip0126

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestUnspecified(t *testing.T) {
	tests := []struct {
		testName  string
		ValueName string
		problems  testutils.Problems
	}{
		{"Valid", "BOOK_FORMAT_UNSPECIFIED", testutils.Problems{}},
		{"InvalidNoPrefix", "UNSPECIFIED", testutils.Problems{{Suggestion: "BOOK_FORMAT_UNSPECIFIED"}}},
		{"InvalidWrongSuffix", "BOOK_FORMAT_UNKNOWN", testutils.Problems{{Suggestion: "BOOK_FORMAT_UNSPECIFIED"}}},
		{"InvalidJustWrong", "UNKNOWN", testutils.Problems{{Suggestion: "BOOK_FORMAT_UNSPECIFIED"}}},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Create the proto with the enum.
			f := testutils.ParseProto3Tmpl(t, `
				enum BookFormat {
					{{.ValueName}} = 0;
					HARDBACK = 1;
					PAPERBACK = 2;
				}
			`, test)

			// Run the lint rule and establish we get the correct problems.
			problems := unspecified.Lint(f)
			if diff := test.problems.SetDescriptor(f.GetEnumTypes()[0].GetValues()[0]).Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Create the proto with the enum.
			f := testutils.ParseProto3Tmpl(t, `
				enum BookFormat {
					option allow_alias = true; 
					HARDBACK = 0;
					{{.ValueName}} = 0;
					PAPERBACK = 2;
				}
			`, test)

			// Run the lint rule and establish we get the correct problems.
			problems := unspecified.Lint(f)
			if diff := test.problems.SetDescriptor(f.GetEnumTypes()[0].GetValues()[0]).Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
