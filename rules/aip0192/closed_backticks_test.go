// Copyright 2023 Google LLC
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

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestClosedBackticks(t *testing.T) {
	for _, test := range []struct {
		name     string
		Comment  string
		problems testutils.Problems
	}{
		{"ValidWordChars", "`foo_bar_baz_1`", nil},
		{"MissingFrontBacktickWordChars", "foo_bar_baz_1`", testutils.Problems{{Suggestion: "`foo_bar_baz_1`"}}},
		{"MissingBackBacktickWordChars", "`foo_bar_baz_1", testutils.Problems{{Suggestion: "`foo_bar_baz_1`"}}},

		{"ValidBetweenWords", "foo `bar` baz", nil},
		{"MissingFrontBacktickBetweenWords", "foo bar` baz", testutils.Problems{{Suggestion: "`bar`"}}},
		{"MissingBackBacktickBetweenWords", "foo `bar baz", testutils.Problems{{Suggestion: "`bar`"}}},

		{"ValidNonWordChars", "`foo:bar/baz->qux.quux[0]()`", nil},
		{"MissingFrontBacktickNonWordChars", "foo:bar/baz->qux.quux[0]()`", testutils.Problems{{Suggestion: "`quux[0]()`"}}},
		{"MissingBackBacktickNonWordChars", "`foo:bar/baz->qux.quux[0]()", testutils.Problems{{Suggestion: "`foo`"}}},

		{"ValidContainsSpace", "`foo + bar`", nil},
		{"MissingFrontBacktickContainsSpace", "foo + bar`", testutils.Problems{{Suggestion: "`bar`"}}},
		{"MissingBackBacktickContainsSpace", "`foo + bar", testutils.Problems{{Suggestion: "`foo`"}}},

		{"ValidNonAscii", "`汉语`", nil},
		{"MissingFrontBacktickNonAscii", "汉语`", testutils.Problems{{Suggestion: "`汉语`"}}},
		{"MissingBackBacktickNonAscii", "`汉语", testutils.Problems{{Suggestion: "`汉语`"}}},

		{"ValidQuotes", "`\"foo\"`", nil},
		{"MissingFrontBacktickQuotes", "\"foo\"`", testutils.Problems{{Suggestion: "`\"foo\"`"}}},
		{"MissingBackBacktickQuotes", "`\"foo\"", testutils.Problems{{Suggestion: "`\"foo\"`"}}},

		{"ValidSeparatedByColons", "foo:`bar:baz`:qux", nil},
		{"MissingFrontBacktickSeparatedByColons", "foo:`bar:baz:qux", testutils.Problems{{Suggestion: "`bar`"}}},
		{"MissingBackBacktickSeparatedByColons", "foo:bar:baz`:qux", testutils.Problems{{Suggestion: "`baz`"}}},

		{"ValidComma", "`name`, a string", nil},
		{"MissingFrontBacktickComma", "name`, a string", testutils.Problems{{Suggestion: "`name`"}}},
		{"MissingBackBacktickComma", "`name, a string", testutils.Problems{{Suggestion: "`name`"}}},

		{"ValidMultipleInlineCode", "`name`: `string`", nil},
		{"MissingFrontBacktickMultipleInlineCode", "name`: string`", testutils.Problems{{Suggestion: "`name`"}, {Suggestion: "`string`"}}},
		{"MissingBackBacktickMultipleInlineCode", "`name: `string", testutils.Problems{{Suggestion: "`name`"}, {Suggestion: "`string`"}}},
		{"MissingFrontAndBackBacktickMultipleInlineCode", "name`: `string", testutils.Problems{{Suggestion: "`name`"}, {Suggestion: "`string`"}}},

		{"ValidContainsColon", "`name: string`", nil},

		{"ValidContainsSeparator", "`.`", nil},
		{"MissingFrontBacktickContainsSeparator", ".`", testutils.Problems{{Suggestion: "``"}}},
		{"MissingBackBacktickContainsSeparator", "`.", testutils.Problems{{Suggestion: "``"}}},

		{"ValidEmptySpace", "``", nil},
		{"ValidEmptySpaces", "`` ``", nil},
		{"UnpairedEmptySpace", "`` `", testutils.Problems{{Suggestion: "``"}}},

		{"IgnoreTripleBacktick", "```", nil},
		{"IgnoreBackticksWithinWords", "a`b`c`d", nil},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
			  // {{.Comment}}
			  message Foo {}
			`, test)
			m := f.GetMessageTypes()[0]
			if diff := test.problems.SetDescriptor(m).Diff(closedBackticks.Lint(f)); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
