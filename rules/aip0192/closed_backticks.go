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
	"regexp"
	"strings"

	"bitbucket.org/creachadair/stringset"
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

type backtickPair int

const (
	unspecified backtickPair = iota
	missingOpening
	missingClosing
	paired
)

// Represents a backtick character in a line.
//
// To narrow which backticks to lint and the likely intended use of a stray
// backtick, we compute whether the backtick can be used as an opening and/or
// closing backtick.
//
// If the character before a backtick is one of the `separators`, the backtick
// is considered opening. Likewise, if the character after a backtick is one of
// the `separators`, the backtick is considered closing.
type backtick struct {
	index   int
	opening bool
	closing bool
	pair    backtickPair
}

var separators = stringset.New(" ", ":", ",", ".", ";", "-")
var reBacktick = regexp.MustCompile("`")

var closedBackticks = &lint.DescriptorRule{
	Name: lint.NewRuleName(192, "closed-backticks"),
	LintDescriptor: func(d desc.Descriptor) []lint.Problem {
		problems := []lint.Problem{}

		for _, comment := range utils.SeparateInternalComments(d.GetSourceInfo().GetLeadingComments()).External {
			for _, line := range strings.Split(comment, "\n") {
				// Pad the line with whitespace, so it's easier to lint
				// backticks at the start and end of the line.
				line = " " + line + " "

				// Find all the backticks in the line and compute whether each
				// one can be used as an opening and/or closing backtick.
				backticks := []backtick{}
				for _, loc := range reBacktick.FindAllStringIndex(line, -1) {
					backticks = append(backticks, backtickAtIndex(line, loc[0]))
				}

				// Compute whether the backticks are paired.
				backticks = computeBacktickPairs(filterUnusableBackticks(backticks))

				// Add a problem for each backtick that is missing a pair.
				for _, backtick := range backticks {
					if backtick.pair != paired {
						problems = append(problems, lint.Problem{
							Message:    "Inline code should be surrounded by backticks.",
							Suggestion: "`" + suggestionWord(line, backtick) + "`",
							Descriptor: d,
						})
					}
				}
			}
		}
		return problems
	},
}

func backtickAtIndex(line string, index int) backtick {
	return backtick{
		index:   index,
		opening: separators.Contains(string(line[index-1])),
		closing: separators.Contains(string(line[index+1])),
	}
}

// Filter out backticks that cannot be used as opening or closing backticks.
func filterUnusableBackticks(backticks []backtick) []backtick {
	list := []backtick{}
	for _, backtick := range backticks {
		if backtick.opening || backtick.closing {
			list = append(list, backtick)
		}
	}
	return list
}

// Compute whether each backtick is paired with an adjacent backtick.
//
// A backtick pair consists of a set of adjacent backticks, where the first is
// opening and the second is closing. Each backtick can be in at most one pair.
//
// This fills in the `pair` field for each backtick. It assumes every backtick
// is either opening, closing, or both.
func computeBacktickPairs(backticks []backtick) []backtick {
	computed := []backtick{}

	prevBacktickOpen := false
	for i, backtick := range backticks {
		backtick.pair = computePairState(backtick, prevBacktickOpen)

		// If this backtick got paired, mark the previous one as paired as well.
		if backtick.pair == paired {
			computed[i-1].pair = paired
		}

		// If paired, this backtick cannot be used as an `opening` for another
		// backtick.
		prevBacktickOpen = false
		if backtick.opening && backtick.pair != paired {
			prevBacktickOpen = true
		}

		computed = append(computed, backtick)
	}
	return computed
}

func computePairState(backtick backtick, prevBacktickOpen bool) backtickPair {
	// If the backtick is only opening, it needs a closing backtick.
	if backtick.opening && !backtick.closing {
		return missingClosing
	}
	// Otherwise, the backtick is closing. It is paired if the last backtick was
	// opening.
	if prevBacktickOpen {
		return paired
	}
	return missingOpening
}

// Extract the word around an unpaired backtick.
//
// Even though inline code can include separators (e.g. colons, spaces), just
// suggest the string of characters up to the first separator character.
func suggestionWord(line string, backtick backtick) string {
	switch backtick.pair {
	case missingOpening:
		return reverseString(firstWordBeforeSeparator(reverseString(line[:backtick.index])))
	case missingClosing:
		return firstWordBeforeSeparator(line[backtick.index+1:])
	default:
		return ""
	}
}

func firstWordBeforeSeparator(s string) string {
	index := strings.IndexFunc(s, func(r rune) bool {
		return separators.Contains(string(r))
	})
	if index != -1 {
		return s[:index]
	}
	return s
}

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
