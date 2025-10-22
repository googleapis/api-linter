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

package aip0192

import (
	"testing"

	"github.com/googleapis/api-linter/v2/rules/internal/testutils"
)

func TestOnlyLeadingComments(t *testing.T) {
	tests := []struct {
		testName string
		Detached string
		Trailing string
		problems testutils.Problems
	}{
		{"ValidNone", "", "", testutils.Problems{}},
		{"ValidInternalDetached", "// (-- detached --)", "", testutils.Problems{}},
		{"ValidInternalTrailing", "", "// (-- trailing --)", testutils.Problems{}},
		{"ValidInternalTrailingNoClose", "", "// (-- trailing", testutils.Problems{}},
		{
			"ValidInternalBoth",
			"// (-- detached --)",
			"// (-- trailing --)",
			testutils.Problems{},
		},
		{
			"ValidInternalBothDouble",
			"// (-- detached --)\n// (-- detached --)",
			"// (-- trailing --)\n// (-- trailing --)",
			testutils.Problems{},
		},
		{"InvalidDetached", "// detached", "", testutils.Problems{{Message: "empty lines"}}},
		{"InvalidTrailing", "", "// trailing", testutils.Problems{{Message: "trailing"}}},
		{
			"InvalidBoth",
			"// detached",
			"// trailing",
			testutils.Problems{{Message: "trailing"}, {Message: "empty lines"}},
		},
		{
			"InvalidPartialDetachedInternalFirst",
			"// (-- detached --) detached",
			"",
			testutils.Problems{{Message: "empty lines"}},
		},
		{
			"InvalidPartialDetachedExternalFirst",
			"// detached (-- detached --)",
			"",
			testutils.Problems{{Message: "empty lines"}},
		},
		{
			"InvalidPartialDetachedSandwich",
			"// (-- detached --) detached (-- detached --)",
			"",
			testutils.Problems{{Message: "empty lines"}},
		},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
				{{.Detached}}

				// The book.
				message Book {
					{{.Trailing}}

					string name = 1;
				}
				
			`, test)
			problems := onlyLeadingComments.Lint(file)
			if diff := test.problems.SetDescriptor(file.Messages().Get(0)).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}
