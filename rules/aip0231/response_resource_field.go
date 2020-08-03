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
	"strings"

	"github.com/gertd/go-pluralize"
	"github.com/googleapis/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
)

// The Batch Get response message should have resource field.
var resourceField = &lint.MessageRule{
	Name:   lint.NewRuleName(231, "response-resource-field"),
	OnlyIf: isBatchGetResponseMessage,
	LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
		// The singular form the resource message name; the first letter capitalized.
		plural := strings.TrimSuffix(strings.TrimPrefix(m.GetName(), "BatchGet"), "Response")
		resourceMsgName := pluralize.NewClient().Singular(plural)

		for _, fieldDesc := range m.GetFields() {
			if msgDesc := fieldDesc.GetMessageType(); msgDesc != nil && msgDesc.GetName() == resourceMsgName {
				if !fieldDesc.IsRepeated() {
					return []lint.Problem{{
						Message:    fmt.Sprintf("The %q type field on Batch Get Response message should be repeated", msgDesc.GetName()),
						Descriptor: fieldDesc,
					}}
				}

				return nil
			}
		}

		// Rule check: Establish that a resource field must be included.
		return []lint.Problem{{
			Message:    fmt.Sprintf("Message %q has no %q type field", m.GetName(), resourceMsgName),
			Descriptor: m,
		}}
	},
}
