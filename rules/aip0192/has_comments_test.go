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

func TestFieldHasComments(t *testing.T) {
	for _, tst := range []struct {
		testName string
		Comment  string
		problems testutils.Problems
	}{
		{"Valid", "This is the title.", nil},
		{"Invalid", "", testutils.Problems{{Message: `Missing comment over`}}},
		{"InvalidInternal", "(-- Internal only comment --)", testutils.Problems{{Message: `Missing comment over`}}},
	} {
		file := testutils.ParseProto3Tmpl(t, `
		// This is a book.
		message Book {
			// The resource name.
			string name = 1;
			// {{ .Comment }}
			string title = 2;
		}
	`, tst)
		problems := tst.problems.SetDescriptor(file.Messages().Get(0).Fields().Get(1))
		if diff := problems.Diff(hasComments.Lint(file)); diff != "" {
			t.Errorf("%s: got(+),want(-):\n%s", tst.testName, diff)
		}
	}
}

func TestMessageHasComments(t *testing.T) {
	for _, tst := range []struct {
		testName string
		Comment  string
		problems testutils.Problems
	}{
		{"Valid", "This is a book.", nil},
		{"Invalid", "", testutils.Problems{{Message: `Missing comment over`}}},
		{"InvalidInternal", "(-- Internal only comment --)", testutils.Problems{{Message: `Missing comment over`}}},
	} {
		file := testutils.ParseProto3Tmpl(t, `
		// {{ .Comment }}
		message Book {
			// The resource name.
			string name = 1;
		}
	`, tst)
		problems := tst.problems.SetDescriptor(file.Messages().Get(0))
		if diff := problems.Diff(hasComments.Lint(file)); diff != "" {
			t.Errorf("%s: got(+),want(-):\n%s", tst.testName, diff)
		}
	}
}
