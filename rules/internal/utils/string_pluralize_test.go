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

package utils

import (
	"testing"
)

func TestPluralize(t *testing.T) {
	tests := []struct {
		name           string
		word           string
		pluralizedWord string
	}{
		{"PluralizeSingularWord", "person", "people"},
		{"PluralizePluralWord", "people", "people"},
		{"PluralizeNonstandardPluralWord", "persons", "people"},
		{"PluralizeNoPluralFormWord", "moose", "moose"},
		{"PluralizePluralLatinWord", "cacti", "cacti"},
		{"PluralizeNonstandardPluralLatinWord", "cactuses", "cacti"},
		{"PluralizePluralCamelCaseWord", "student_profiles", "student_profiles"},
		{"PluralizeSingularCamelCaseWord", "student_profile", "student_profiles"},
	}
	for _, test := range tests {
		t.Run(test.word, func(t *testing.T) {
			if got := ToPlural(test.word); got != test.pluralizedWord {
				t.Errorf("Plural(%s) got %s, but want %s", test.word, got, test.pluralizedWord)
			}
		})
	}
}
