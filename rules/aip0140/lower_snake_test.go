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

	"github.com/googleapis/api-linter/rules/internal/testutils"
	"github.com/jhump/protoreflect/desc/builder"
)

func TestLowerSnake(t *testing.T) {
	// Define permutations.
	tests := []struct {
		testName  string
		fieldName string
		problems  testutils.Problems
	}{
		{"ValidOneWord", "rated", testutils.Problems{}},
		{"ValidTwoWords", "has_rating", testutils.Problems{}},
		{"InvalidCamel", "hasRating", testutils.Problems{{Suggestion: "has_rating"}}},
		{"InvalidConstantOneWord", "RATED", testutils.Problems{{Suggestion: "rated"}}},
		{"InvalidConstantTwoWords", "HAS_RATING", testutils.Problems{{Suggestion: "has_rating"}}},
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

			// Run the lint rule and verify that we got the expected set
			// of problems.
			problems := lowerSnake.Lint(message.GetFile())
			if diff := test.problems.SetDescriptor(message.GetFields()[0]).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}
