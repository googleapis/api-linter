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

package aip0131

import (
	"fmt"

	"github.com/googleapis/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"
)

// The Get standard method should only have expected fields.
var standardFields = lint.Rule{
	Name:        lint.NewRuleName("core", "0131", "request-message", "name-field"),
	Description: "The Get standard method must include expected fields.",
	URI:         "https://aip.dev/131#request-message",
	LintMessage: func(m *desc.MessageDescriptor) (problems []lint.Problem) {
		// We only care about Get methods for the purpose of this rule;
		// ignore everything else.
		if !isGetRequestMessage(m) {
			return problems
		}

		// Rule check: Establish that a name field is present.
		name := m.FindFieldByName("name")
		if name == nil {
			problems = append(problems, lint.Problem{
				Message:    fmt.Sprintf("Method %q has no `name` field", m.GetName()),
				Descriptor: m,
			})
			return problems
		}

		// Rule check: Ensure that the name field is the correct type.
		if name.GetType() != builder.FieldTypeString().GetType() {
			problems = append(problems, lint.Problem{
				Message:    "`name` field on Get RPCs should be a string",
				Descriptor: name,
			})
		}

		return problems
	},
}

// Get methods should not have unrecognized fields.
var unknownFields = lint.Rule{
	Name:        lint.NewRuleName("core", "0131", "request-message", "unknown-fields"),
	Description: "Get RPCs must not contain unexpected fields.",
	URI:         "https://aip.dev/131#request-message",
	LintMessage: func(m *desc.MessageDescriptor) (problems []lint.Problem) {
		// We only care about Get methods for the purpose of this rule;
		// ignore everything else.
		if !isGetRequestMessage(m) {
			return
		}

		// Rule check: Establish that there are no unexpected fields.
		allowedFields := map[string]struct{}{
			"name":      {}, // AIP-131
			"read_mask": {}, // AIP-157
			"view":      {}, // AIP-157
		}
		for _, field := range m.GetFields() {
			if _, ok := allowedFields[string(field.GetName())]; !ok {
				problems = append(problems, lint.Problem{
					Message: fmt.Sprintf(
						"Get RPCs must only contain fields explicitly described in AIPs, not %q.",
						string(field.GetName()),
					),
					Descriptor: field,
				})
			}
		}

		return problems
	},
}
