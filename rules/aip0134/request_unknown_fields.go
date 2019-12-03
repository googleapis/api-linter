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

package aip0134

import (
	"fmt"

	"bitbucket.org/creachadair/stringset"
	"github.com/googleapis/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
)

// Update methods should not have unrecognized fields.
var unknownFields = &lint.MessageRule{
	Name:   lint.NewRuleName(134, "request-unknown-fields"),
	OnlyIf: isUpdateRequestMessage,
	LintMessage: func(m *desc.MessageDescriptor) (problems []lint.Problem) {
		resource := extractResource(m.GetName())
		// Rule check: Establish that there are no unexpected fields.
		allowedFields := stringset.New(
			fieldNameFromResource(resource), // AIP-134
			"request_id",                    // AIP-155
			"update_mask",                   // AIP-134
			"validate_only",                 // AIP needed
		)
		for _, field := range m.GetFields() {
			if !allowedFields.Contains(field.GetName()) {
				problems = append(problems, lint.Problem{
					Message: fmt.Sprintf(
						"Unexpected field: Update RPCs must only contain fields explicitly described in AIPs, not %q.",
						field.GetName(),
					),
					Descriptor: field,
				})
			}
		}

		return
	},
}
