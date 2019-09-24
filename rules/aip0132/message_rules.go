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

package aip0132

import (
	"fmt"

	"github.com/googleapis/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"
)

// The List standard method should contain a parent field.
var standardFields = &lint.MessageRule{
	Name: lint.NewRuleName("core", "0132", "request-message", "parent-field"),
	URI:  "https://aip.dev/132#request-message",
	LintMessage: func(m *desc.MessageDescriptor) lint.Problems {
		// We only care about List- methods for the purpose of this rule;
		// ignore everything else.
		if !isListRequestMessage(m) {
			return nil
		}

		// Rule check: Establish that a `parent` field is present.
		parentField := m.FindFieldByName("parent")
		if parentField == nil {
			return lint.Problems{{
				Message:    fmt.Sprintf("Message %q has no `parent` field", m.GetName()),
				Descriptor: m,
			}}
		}

		// Rule check: Establish that the parent field is a string.
		if parentField.GetType() != builder.FieldTypeString().GetType() {
			return lint.Problems{{
				Message:    "`parent` field on List RPCs must be a string",
				Descriptor: parentField,
			}}
		}

		return nil
	},
}

// List methods should not have unrecognized fields.
var unknownFields = &lint.MessageRule{
	Name: lint.NewRuleName("core", "0132", "request-message", "unknown-fields"),
	URI:  "https://aip.dev/132#request-message",
	LintMessage: func(m *desc.MessageDescriptor) (problems lint.Problems) {
		// We only care about List- methods for the purpose of this rule;
		// ignore everything else.
		if !isListRequestMessage(m) {
			return
		}

		// Rule check: Establish that there are no unexpected fields.
		//
		// Additionally, we type check only the fields defined in AIP-132,
		// but leave fields defined elsewhere to be type checked by those linter
		// rules.
		allowedFields := map[string]*builder.FieldType{
			"parent":       builder.FieldTypeString(), // AIP-132
			"page_size":    nil,                       // AIP-158
			"page_token":   nil,                       // AIP-158
			"filter":       builder.FieldTypeString(), // AIP-132
			"order_by":     builder.FieldTypeString(), // AIP-132
			"group_by":     builder.FieldTypeString(), // Nowhere yet, but permitted
			"show_deleted": nil,                       // AIP-135
			"read_mask":    nil,                       // AIP-157
			"view":         nil,                       // AIP-157
		}
		for _, field := range m.GetFields() {
			fieldType, allowed := allowedFields[field.GetName()]
			if !allowed {
				problems = append(problems, lint.Problem{
					Message:    "List RPCs should only contain fields explicitly described in AIPs.",
					Descriptor: field,
				})
			}
			if fieldType != nil && field.GetType() != fieldType.GetType() {
				problems = append(problems, lint.Problem{
					Message: fmt.Sprintf(
						"Field %q is the wrong type; expected %q.",
						field.GetName(), fieldType.GetTypeName()),
					Descriptor: field,
				})
			}
		}

		return
	},
}
