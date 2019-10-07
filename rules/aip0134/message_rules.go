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
	"strings"

	"github.com/googleapis/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
)

func extractResource(reqName string) string {
	// Strips "Update" from the beginning and "Request" from the end.
	return reqName[6:len(reqName)-7]
}

// The Update standard method should only have expected fields.
var standardFields = &lint.MessageRule{
	Name:   lint.NewRuleName("core", "0134", "request-message", "name-field"),
	URI:    "https://aip.dev/134#request-message",
	OnlyIf: isUpdateRequestMessage,
	LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
		r := extractResource(m.GetName())
		fields := []struct {
			name string
			typ  string
		}{
			{strings.ToLower(r), r},
			{"update_mask", "google.protobuf.FieldMask"},
		}
		// Rule check: Establish that expected fields are present.
		for _, f := range fields {
			// Rule check: Establish that the field is present.
			field := m.FindFieldByName(f.name)
			if field == nil {
				return []lint.Problem{{
					Message:    fmt.Sprintf("Method %s has no `%s` field", m.GetName(), f.name),
					Descriptor: m,
				}}
			}

			// Rule check: Ensure that the field is the correct type.
			msgDesc := field.GetMessageType()
			if msgDesc == nil || msgDesc.GetFullyQualifiedName() != f.typ {
				return []lint.Problem{{
					Message:    fmt.Sprintf("`%s` field on Update RPCs should be of type `%s`", f.name, f.typ),
					Descriptor: field,
				}}
			}
		}

		return nil
	},
}

// Update methods should not have unrecognized fields.
var unknownFields = &lint.MessageRule{
	Name:   lint.NewRuleName("core", "0134", "request-message", "unknown-fields"),
	URI:    "https://aip.dev/134#request-message",
	OnlyIf: isUpdateRequestMessage,
	LintMessage: func(m *desc.MessageDescriptor) (problems []lint.Problem) {
		resource := extractResource(m.GetName());
		// Rule check: Establish that there are no unexpected fields.
		allowedFields := map[string]struct{}{
			strings.ToLower(resource): {}, // AIP-134
			"update_mask":             {}, // AIP-134
		}
		for _, field := range m.GetFields() {
			if _, ok := allowedFields[string(field.GetName())]; !ok {
				problems = append(problems, lint.Problem{
					Message: fmt.Sprintf(
						"Unexpected field: Update RPCs must only contain fields explicitly described in AIPs, not %q.",
						string(field.GetName()),
					),
					Descriptor: field,
				})
			}
		}

		return
	},
}
