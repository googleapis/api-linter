package aip0126

import (
	"testing"

	"github.com/jhump/protoreflect/desc/builder"
)

func TestUpperSnake(t *testing.T) {
	// Define permutations.
	tests := []struct {
		name             string
		enumValues       []string
		wantProblemCount int
		wantSuggestions  []string
	}{
		{"ValidOneWord", []string{"ONE"}, 0, nil},
		{"ValidTwoWords", []string{"ONE_TWO"}, 0, nil},
		{"InvalidOneWord", []string{"one"}, 1, []string{"ONE"}},
		{"InvalidTwoWordsCamel", []string{"oneTwo"}, 1, []string{"ONE_TWO"}},
		{"InvalidTwoWordsLowerSnake", []string{"one_two"}, 1, []string{"ONE_TWO"}},
		{"OneProblem", []string{"one_two", "THREE_FOUR"}, 1, []string{"ONE_TWO"}},
		{"TwoProblems", []string{"one_two", "three_four"}, 2, []string{"ONE_TWO", "THREE_FOUR"}},
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

			problems := enumValueUpperSnakeCase.LintEnum(enum)
			if got, want := len(problems), test.wantProblemCount; got != want {
				t.Errorf("rule enumValueUpperSnakeCase got %d problems, but wanted %d", got, want)
			}

			for i, problem := range problems {
				if got, want := problem.Suggestion, test.wantSuggestions[i]; got != want {
					t.Errorf("rule enumValueUpperSnakeCase got suggestion %q, but wanted %q", got, want)
				}
			}
		})
	}
}
