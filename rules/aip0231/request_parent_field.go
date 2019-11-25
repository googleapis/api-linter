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

package aip0231

import (
	"fmt"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"
)

// The Batch Get request message should have parent field.
var parentField = &descrule.MessageRule{
	RuleName:   lint.NewRuleName(231, "request-parent-field"),
	OnlyIf: isBatchGetRequestMessage,
	LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
		// Rule check: Establish that a `parent` field is present.
		parentField := m.FindFieldByName("parent")
		if parentField == nil {
			return []lint.Problem{{
				Message:    fmt.Sprintf("Message %q has no `parent` field", m.GetName()),
				Descriptor: m,
			}}
		}

		// Rule check: Establish that the parent field is a string.
		if parentField.GetType() != builder.FieldTypeString().GetType() {
			return []lint.Problem{{
				Message:    "`parent` field on create request message must be a string",
				Suggestion: "string",
				Descriptor: parentField,
				Location:   locations.FieldType(parentField),
			}}
		}

		return nil
	},
}
