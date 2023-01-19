// Copyright 2020 Google LLC
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

package aip0216

import (
	"fmt"
	"strings"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/jhump/protoreflect/desc"
)

var valueSynonyms = &lint.EnumValueRule{
	Name: lint.NewRuleName(216, "value-synonyms"),
	OnlyIf: func(v *desc.EnumValueDescriptor) bool {
		return strings.HasSuffix(v.GetEnum().GetName(), "State")
	},
	LintEnumValue: func(v *desc.EnumValueDescriptor) []lint.Problem {
		for bad, good := range map[string]string{
			"CANCELED":   "CANCELLED",
			"CANCELING":  "CANCELLING",
			"FAIL":       "FAILED",
			"FAILURE":    "FAILED",
			"READY":      "ACTIVE",
			"SUCCESS":    "SUCCEEDED",
			"SUCCESSFUL": "SUCCEEDED",
		} {
			if v.GetName() == bad {
				return []lint.Problem{{
					Message:    fmt.Sprintf("Prefer %q over %q for state names.", good, bad),
					Suggestion: good,
					Descriptor: v,
					Location:   locations.DescriptorName(v),
				}}
			}
		}
		return nil
	},
}
