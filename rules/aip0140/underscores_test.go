// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aip0140

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestUnderscores(t *testing.T) {
	for _, test := range []struct {
		name     string
		Field    string
		problems testutils.Problems
	}{
		{"Valid", "foo_bar", nil},
		{"InvalidLeading", "_foo_bar", testutils.Problems{{Suggestion: "foo_bar"}}},
		{"InvalidTrailing", "foo_bar_", testutils.Problems{{Suggestion: "foo_bar"}}},
		{"InvalidDoubleLeading", "__foo_bar", testutils.Problems{{Suggestion: "foo_bar"}}},
		{"InvalidDoubleTrailing", "foo_bar__", testutils.Problems{{Suggestion: "foo_bar"}}},
		{"InvalidAdjacent", "foo__bar", testutils.Problems{{Suggestion: "foo_bar"}}},
		{"InvalidDoubleAdjacent", "foo___bar", testutils.Problems{{Suggestion: "foo_bar"}}},
		{"InvalidMultiAdjacent", "foo__bar__baz", testutils.Problems{{Suggestion: "foo_bar_baz"}}},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				message Foo {
					string {{.Field}} = 1;
				}
			`, test)
			field := f.GetMessageTypes()[0].GetFields()[0]
			if diff := test.problems.SetDescriptor(field).Diff(underscores.Lint(f)); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
