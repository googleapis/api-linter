package aip0126

import (
	"testing"

	"github.com/commure/api-linter/rules/internal/testutils"
	"github.com/jhump/protoreflect/desc/builder"
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
			enumBuilder := builder.NewEnum("Number")
			for _, enumValue := range test.enumValues {
				enumBuilder.AddValue(builder.NewEnumValue(enumValue))
			}

			enum, err := enumBuilder.Build()
			if err != nil {
				t.Fatalf("Could not build an enum with values %v", test.enumValues)
			}

			// If this test expects problems, they correspond 1:1 with the
			// enum value order.
			for i := range test.problems {
				test.problems[i].Descriptor = enum.GetValues()[i]
			}

			// Run the lint rule, and establish that we got the correct problems.
			problems := enumValueUpperSnakeCase.Lint(enum.GetFile())
			if diff := test.problems.Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
