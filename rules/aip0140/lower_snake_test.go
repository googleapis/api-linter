// Copyright 2019 Google LLC
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

	"github.com/jhump/protoreflect/desc/builder"
)

func TestLowerSnake(t *testing.T) {
	// Define permutations.
	tests := []struct {
		testName     string
		fieldName    string
		problemCount int
		errPrefix    string
	}{
		{"ValidOneWord", "rated", 0, "False positive"},
		{"ValidTwoWords", "has_rating", 0, "False positive"},
		{"InvalidCamel", "hasRating", 1, "False negative"},
		{"InvalidConstantOneWord", "RATED", 1, "False negative"},
		{"InvalidConstantTwoWords", "HAS_RATING", 1, "False negative"},
	}

	// Test each permutation.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Create the given field.
			message, err := builder.NewMessage("Book").AddField(
				builder.NewField(test.fieldName, builder.FieldTypeBool()),
			).Build()
			if err != nil {
				t.Fatalf("Could not build `%s` field.", test.fieldName)
			}
			field := message.FindFieldByName(test.fieldName)

			// Run the lint rule and verify that we got the expected set
			// of problems.
			if problems := lowerSnake.LintField(field); len(problems) != test.problemCount {
				t.Errorf("%s on rule %s: %#v", test.errPrefix, lowerSnake.Name, problems)
			}
		})
	}
}
