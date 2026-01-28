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

package utils

// ToUpperCamelCase returns the UpperCamelCase of a string, including removing
// delimiters (_,-,., ) and using them to denote a new word.
func ToUpperCamelCase(s string) string {
	return toCamelCase(s, true, false)
}

// ToLowerCamelCase returns the lowerCamelCase of a string, including removing
// delimiters (_,-,., ) and using them to denote a new word.
func ToLowerCamelCase(s string) string {
	return toCamelCase(s, false, true)
}

func toCamelCase(s string, makeNextUpper bool, makeNextLower bool) string {
	asLower := make([]rune, 0, len(s))
	for _, r := range s {
		if isLower(r) {
			if makeNextUpper {
				r = r & '_' // make uppercase
			}
			asLower = append(asLower, r)
		} else if isUpper(r) {
			if makeNextLower {
				r = r | ' ' // make lowercase
			}
			asLower = append(asLower, r)
		} else if isNumber(r) {
			asLower = append(asLower, r)
		}
		makeNextUpper = false
		makeNextLower = false

		if r == '-' || r == '_' || r == ' ' || r == '.' {
			// handle snake case scenarios, which generally indicates
			// a delimited word.
			makeNextUpper = true
		}
	}
	return string(asLower)
}

func isUpper(r rune) bool {
	return ('A' <= r && r <= 'Z')
}

func isNumber(r rune) bool {
	return ('0' <= r && r <= '9')
}

func isLower(r rune) bool {
	return ('a' <= r && r <= 'z')
}
