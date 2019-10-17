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

// Package aip0143 contains rules defined in https://aip.dev/143.
package aip0143

import (
	"github.com/googleapis/api-linter/lint"
)

// AddRules adds all of the AIP-141 rules to the provided registry.
func AddRules(r lint.RuleRegistry) {
	r.Register(
		fieldNames,
		fieldTypes,
	)
}

var standardFields = make(map[string]struct{})

func init() {
	for _, n := range []string{
		"mime_type",
		"country_code",
		"currency_code",
		"language_code",
		"time_zone",
	} {
		standardFields[n] = struct{}{}
	}
}
