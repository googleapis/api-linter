package aip0126

import (
	"fmt"
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestUpperSnake(t *testing.T) {
	// Define permutations.
	tests := []struct {
		name       string
		enumValues []string
		problems   testutils.Problems
	}{
		{"ValidOneWord", []string{"ONE"}, testutils.Problems{}},
		{"ValidTwoWords", []string{"ONE_TWO"}, testutils.Problems{}},
		{"ValidTwoWordsTrailingNumber", []string{"ONE_TWO2"}, testutils.Problems{}},
		{"ValidThreeWordsTrailingNumber", []string{"ONE_TWO2_THREE"}, testutils.Problems{}},
		{"ValidThreeWordsTrailingLeadingNumber", []string{"ONE_TWO2_3THREE"}, testutils.Problems{}},
		{"ValidThreeWordsIsolatedNumber", []string{"ONE_TWO_3"}, testutils.Problems{}},
		{"InvalidOneWord", []string{"one"}, testutils.Problems{{Suggestion: "ONE"}}},
		{"InvalidTwoWordsCamel", []string{"oneTwo"}, testutils.Problems{{Suggestion: "ONE_TWO"}}},
		{"InvalidTwoWordsLowerSnake", []string{"one_two"}, testutils.Problems{{Suggestion: "ONE_TWO"}}},
		{"OneProblem", []string{"one_two", "THREE_FOUR"}, testutils.Problems{{Suggestion: "ONE_TWO"}}},
		{
			"TwoProblems",
			[]string{"one_two", "three_four"},
			testutils.Problems{{Suggestion: "ONE_TWO"}, {Suggestion: "THREE_FOUR"}},
		},
	}

	// Test each permutation.
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var enumValues string
			for i, v := range test.enumValues {
				enumValues += fmt.Sprintf("%s = %d;\n", v, i)
			}
			f := testutils.ParseProto3Tmpl(t, `
				enum Number {
					{{.EnumValues}}
				}
			`, struct{ EnumValues string }{enumValues})
			enum := f.Enums().Get(0)

			// If this test expects problems, they correspond 1:1 with the
			// enum value order.
			for i := range test.problems {
				test.problems[i].Descriptor = enum.Values().Get(i)
			}

			// Run the lint rule, and establish that we got the correct problems.
			problems := enumValueUpperSnakeCase.Lint(f)
			if diff := test.problems.Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}
