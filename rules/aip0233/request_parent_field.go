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

package aip0233

import (
	"fmt"
	"strings"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"
	"github.com/stoewer/go-strcase"
)

// The Batch Create request message should have parent field.
var requestParentField = &lint.MessageRule{
	Name: lint.NewRuleName(233, "request-parent-field"),
	OnlyIf: func(m *desc.MessageDescriptor) bool {
		// Sanity check: If the resource has a pattern, and that pattern
		// contains only one variable, then a parent field is not expected.
		//
		// In order to parse out the pattern, we get the resource message
		// from the response, then get the resource annotation from that,
		// and then inspect the pattern there (oy!).
		plural := strings.TrimPrefix(strings.TrimSuffix(m.GetName(), "Request"), "BatchCreate")
		if resp := m.GetFile().FindMessage(fmt.Sprintf("BatchCreate%sResponse", plural)); resp != nil {
			if paged := resp.FindFieldByName(strcase.SnakeCase(plural)); paged != nil {
				if resource := utils.GetResource(paged.GetMessageType()); resource != nil {
					for _, pattern := range resource.GetPattern() {
						if strings.Count(pattern, "{") == 0 {
							return false
						}
					}
				}
			}
		}

		return isBatchCreateRequestMessage(m)
	},
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
				Message:    "`parent` field on Batch Create request message must be a string.",
				Descriptor: parentField,
				Location:   locations.FieldType(parentField),
				Suggestion: "string",
			}}
		}

		return nil
	},
}
