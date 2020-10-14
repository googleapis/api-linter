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

package aip0164

import (
	"fmt"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"
)

var requestNameField = &lint.MessageRule{
	Name: lint.NewRuleName(164, "request-name-field"),
	OnlyIf: func(m *desc.MessageDescriptor) bool {
		return isUndeleteRequestMessage(m)
	},
	LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
		// Rule check: Establish that a `name` field is present.
		nameField := m.FindFieldByName("name")
		if nameField == nil {
			return []lint.Problem{{
				Message:    fmt.Sprintf("Message %q has no `name` field.", m.GetName()),
				Descriptor: m,
			}}
		}

		// Rule check: Establish that the `name` field is a string.
		if nameField.GetType() != builder.FieldTypeString().GetType() {
			return []lint.Problem{{
				Message:    "`name` field on Undelete request message must be a string.",
				Descriptor: nameField,
				Location:   locations.FieldType(nameField),
				Suggestion: "string",
			}}
		}

		return nil
	},
}
