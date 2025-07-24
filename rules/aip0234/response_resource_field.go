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

package aip0234

import (
	"fmt"
	"strings"

	"github.com/gertd/go-pluralize"
	"github.com/googleapis/api-linter/lint"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// The Batch Update response message should have resource field.
var responseResourceField = &lint.MessageRule{
	Name:   lint.NewRuleName(234, "response-resource-field"),
	OnlyIf: isBatchUpdateResponseMessage,
	LintMessage: func(m protoreflect.MessageDescriptor) []lint.Problem {
		// the singular form the resource name, the first letter is Capitalized.
		// Note: Retrieve the resource name from the the batch update response,
		// for example: "BatchUpdateBooksResponse" -> "Books"
		resourceMsgName := pluralize.NewClient().Singular(strings.TrimPrefix(strings.TrimSuffix(m.Name(), "Response"), "BatchUpdate"))

		for _, fieldDesc := range m.Fields() {
			if msgDesc := fieldDesc.GetMessageType(); msgDesc != nil && msgDesc.Name() == resourceMsgName {
				if !fieldDesc.IsRepeated() {
					return []lint.Problem{{
						Message:    fmt.Sprintf("The %q type field on Batch Update Response message should be repeated", msgDesc.Name()),
						Descriptor: fieldDesc,
					}}
				}

				return nil
			}
		}

		// Rule check: Establish that a resource field must be included.
		return []lint.Problem{{
			Message:    fmt.Sprintf("Message %q has no %q type field", m.Name(), resourceMsgName),
			Descriptor: m,
		}}
	},
}
