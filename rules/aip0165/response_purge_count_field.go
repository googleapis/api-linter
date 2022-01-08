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

	"github.com/commure/api-linter/lint"
	"github.com/commure/api-linter/locations"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"
)

// The Purge response message should have purge_count field.
var responsePurgeCountField = &lint.MessageRule{
	Name:   lint.NewRuleName(165, "response-purge-count-field"),
	OnlyIf: isPurgeResponseMessage,
	LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
		// Rule check: Establish that a `purge_count` field is present.
		field := m.FindFieldByName("purge_count")
		if field == nil {
			return []lint.Problem{{
				Message:    fmt.Sprintf("Message %q has no `purge_count` field.", m.GetName()),
				Descriptor: m,
			}}
		}

		// Rule check: Establish that the purge_count field is a singular int32.
		if field.GetType() != builder.FieldTypeInt32().GetType() || field.IsRepeated() {
			return []lint.Problem{{
				Message:    "`purge_count` field on Purge response message must be a singular int32.",
				Descriptor: field,
				Location:   locations.FieldType(field),
				Suggestion: "int32",
			}}
		}

		return nil
	},
}
