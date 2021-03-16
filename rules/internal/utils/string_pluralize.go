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

package utils

import (
	"github.com/gertd/go-pluralize"
)

var pluralizeClient = pluralize.NewClient()

// ToPlural converts a string to its plural form.
func ToPlural(s string) string {
	// Need to convert name to singular first to support none standard case such as persons, cactuses.
	// persons -> person -> people

	return pluralizeClient.Plural(pluralizeClient.Singular(s))
}

// ToSingular converts a string to its singular form.
func ToSingular(s string) string {
	return pluralizeClient.Singular(s)
}
