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
)

func TestLowerSnake(t *testing.T) {
	// Define permutations.
	tests := []struct {
		testName  string
		FieldName string
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
			f := testutils.ParseProto3Tmpl(t, `
				message Book {
					bool {{.FieldName}} = 1;
				}
			`, test)
			field := f.Messages().Get(0).Fields().Get(0)

			// Run the lint rule and verify that we got the expected set
			// of problems.
			problems := lowerSnake.Lint(f)
			if diff := test.problems.SetDescriptor(field).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}