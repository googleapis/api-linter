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

var fieldNameToCommonType = map[string]string{
	"duration": "google.protobuf.Duration",

	"color":  "google.type.Color",
	"colour": "google.type.Color",

	"dollars": "google.type.Money",
	"euros":   "google.type.Money",
	"money":   "google.type.Money",
	"pounds":  "google.type.Money",
	"yen":     "google.type.Money",
	"yuan":    "google.type.Money",

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
		// Flag this field if it has a name in `fieldNameToCommonType` but
		// doesn't have its corresponding type.
		if messageType, ok := fieldNameToCommonType[f.GetName()]; ok {
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
