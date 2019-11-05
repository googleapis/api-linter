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

	"github.com/gertd/go-pluralize"
	"github.com/googleapis/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
)

// The Batch Get response message should have resource field.
var resourceField = &lint.MessageRule{
	Name:   lint.NewRuleName("core", "0231", "response-resource-field"),
	OnlyIf: isBatchGetResponseMessage,
	LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
		// the singular form the resource name, the first letter is Capitalized.
		// Note: Use m.GetName()[8 : len(m.GetName())-9] to retrieve the resource
		// name from the the batch get response, for example:
		// "BatchGetBooksResponse" -> "Books"
		resourceMsgName := pluralize.NewClient().Singular(m.GetName()[8 : len(m.GetName())-9])

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
