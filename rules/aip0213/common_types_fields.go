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
	"strings"

	"github.com/googleapis/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
)

var fieldNameToCommonType = map[string]string{
	"delay":    "google.protobuf.Duration",
	"duration": "google.protobuf.Duration",
	"timeout":  "google.protobuf.Duration",

	"color":  "google.type.Color",
	"colour": "google.type.Color",

	"balance": "google.type.Money",
	"cost":    "google.type.Money",
	"fare":    "google.type.Money",
	"fee":     "google.type.Money",
	"price":   "google.type.Money",
	"spend":   "google.type.Money",

	"mobile_number": "google.type.PhoneNumber",
	"phone":         "google.type.PhoneNumber",
	"phone_number":  "google.type.PhoneNumber",

	"clock_time":  "google.type.TimeOfDay",
	"time_of_day": "google.type.TimeOfDay",

	"timezone":  "google.type.TimeZone",
	"time_zone": "google.type.TimeZone",
}

var commonTypesFields = &lint.FieldRule{
	Name: lint.NewRuleName(213, "common-types-fields"),
	LintField: func(f *desc.FieldDescriptor) []lint.Problem {
		// Flag this field if it ends with a name in `fieldNameToCommonType` but
		// doesn't have its corresponding type.
		commonFieldName := commonFieldNameSuffix(f.GetName())
		if messageType, ok := fieldNameToCommonType[commonFieldName]; ok {
			if f.GetMessageType() == nil || f.GetMessageType().GetFullyQualifiedName() != messageType {
				return []lint.Problem{{
					Message:    fmt.Sprintf("Consider using the common type %q.", messageType),
					Descriptor: f,
				}}
			}
		}
		return nil
	},
}

// Returns the common field name if `fieldName` ends in one.
func commonFieldNameSuffix(fieldName string) string {
	fieldParts := strings.Split(fieldName, "_")
	for name := range fieldNameToCommonType {
		nameParts := strings.Split(name, "_")
		// Check if `fieldName` ends with the same underscore-separated terms as
		// `name`.
		if strListContains(reverseStrList(fieldParts), reverseStrList(nameParts)) {
			return name
		}
	}
	return ""
}

// Returns true if list `s` contains list `subList` in order.
func strListContains(s []string, subList []string) bool {
	if len(subList) > len(s) {
		return false
	}
	for i := 0; i < len(subList); i++ {
		if s[i] != subList[i] {
			return false
		}
	}
	return true
}

// Returns list `s` in reverse order.
func reverseStrList(s []string) []string {
	out := make([]string, len(s))
	copy(out, s)

	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}
	return out
}
