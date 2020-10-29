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

package aip0165

import (
	"fmt"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"
)

// The Purge request message should have filter field.
var requestFilterField = &lint.MessageRule{
	Name:   lint.NewRuleName(165, "request-filter-field"),
	OnlyIf: isPurgeRequestMessage,
	LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
		// Rule check: Establish that a `filter` field is present.
		filterField := m.FindFieldByName("filter")
		if filterField == nil {
			return []lint.Problem{{
				Message:    fmt.Sprintf("Message %q has no `filter` field.", m.GetName()),
				Descriptor: m,
			}}
		}

		// Rule check: Establish that the filter field is a string.
		if filterField.GetType() != builder.FieldTypeString().GetType() || filterField.IsRepeated() {
			return []lint.Problem{{
				Message:    "`filter` field on Purge request message must be a singular string.",
				Descriptor: filterField,
				Location:   locations.FieldType(filterField),
				Suggestion: "string",
			}}
		}

		return nil
	},
}
