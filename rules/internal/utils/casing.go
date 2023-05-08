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

import "regexp"

var (
	lowerCamelCaseRegex = regexp.MustCompile(`^\p{Ll}[A-Za-z0-9]+$`)
	upperCamelCaseRegex = regexp.MustCompile(`^\p{Lu}[A-Za-z0-9]+$`)
)

// IsLowerCamelCase returns true if the string is lowerCamelCase
//
// lowerCamelCase is defined as:
// - starts with a lowercase letter
// - only includes alphanumeric characters
//
// this is stricter than strcase.LowerCamelCase, which
// allows dashes and underscores.
func IsLowerCamelCase(s string) bool {
	return lowerCamelCaseRegex.MatchString(s)
}

// IsUpperCamelCase returns true if the string is upperCamelCase
//
// upperCamelCase is defined as:
// - starts with an uppercase letter
// - only includes alphanumeric characters
//
// this is stricter than strcase.UpperCamelCase, which
// allows dashes and underscores.
func IsUpperCamelCase(s string) bool {
	return upperCamelCaseRegex.MatchString(s)
}
