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

package aip0213

import (
	"fmt"

	"github.com/googleapis/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
)

// A map from each common type to its constituent fields.
var commonTypesToFields = map[string][]string{
	"google.type.Color":      {"red", "green", "blue"},
	"google.type.Date":       {"year", "month", "day"},
	"google.type.Fraction":   {"numerator", "denominator"},
	"google.type.Interval":   {"start_time", "end_time"},
	"google.type.LatLng":     {"latitude", "longitude"},
	"google.type.TimeOfDay":  {"hours", "minutes", "seconds"},
	"google.type.Quaternion": {"x", "y", "z", "w"},
}

var commonTypesMessages = &lint.MessageRule{
	Name: lint.NewRuleName(213, "common-types-messages"),
	LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
		problems := []lint.Problem{}

		// If a message has all the value fields, it should consider using the
		// key type.
		for commonType, fields := range commonTypesToFields {
			if messageContainsAllFields(m, fields) {
				problems = append(problems, lint.Problem{
					Message:    fmt.Sprintf("Message contains fields %v. Consider using the common type %q.", fields, commonType),
					Descriptor: m,
				})
			}
		}
		return problems
	},
}

func messageContainsAllFields(m *desc.MessageDescriptor, fieldNames []string) bool {
	for _, fieldName := range fieldNames {
		if m.FindFieldByName(fieldName) == nil {
			return false
		}
	}
	return true
}
