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
	"strings"
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestNoHTML(t *testing.T) {
	rawHTMLProblem := testutils.Problems{{Message: "must not include raw HTML"}}
	noProblems := testutils.Problems{}

	for _, test := range []struct {
		name     string
		comment  string
		problems testutils.Problems
	}{
		{"Valid", "It is **great!**", noProblems},
		{"ValidMath", "x < 10", noProblems},
		{"ValidMoreMath", "x < 10 > y", noProblems},
		{"ValidAngleBracketPlaceholders", "Format: http://<server>/<path>", noProblems},
		{"ValidComplexTag", `<img src="https://foo.bar/mickey" />`, noProblems},
		{"InvalidBold", "It is <b>great!</b>", rawHTMLProblem},
		{"InvalidCode", "This is <code>code font</code>.", rawHTMLProblem},
		{"InvalidBreak", "This spans<br />two lines.", rawHTMLProblem},
		{"IntentionalFalseNegativeInnerSpace", "Something < b > bold < /b >", noProblems},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3String(t, strings.ReplaceAll(`
				// A foo. {{.Comment}}
				message Foo {}
			`, "{{.Comment}}", test.comment))
			message := f.GetMessageTypes()[0]
			if diff := test.problems.SetDescriptor(message).Diff(noHTML.Lint(f)); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
