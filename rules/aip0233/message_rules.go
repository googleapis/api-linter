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

	"github.com/gertd/go-pluralize"
	"github.com/googleapis/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"
)

// The Batch Create request message should have parent field.
var parentField = &lint.MessageRule{
	Name:   lint.NewRuleName("core", "0233", "request-message", "parent-field"),
	URI:    "https://aip.dev/233#request-message",
	OnlyIf: isBatchCreateRequestMessage,
	LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
		// Rule check: Establish that a `parent` field is present.
		parentField := m.FindFieldByName("parent")
		if parentField == nil {
			return []lint.Problem{{
				Message:    fmt.Sprintf(`Message %q has no "parent" field`, m.GetName()),
				Descriptor: m,
			}}
		}

		// Rule check: Establish that the parent field is a string.
		if parentField.GetType() != builder.FieldTypeString().GetType() {
			return []lint.Problem{{
				Message:    `"parent" field on create request message must be a string`,
				Descriptor: parentField,
			}}
		}

		return nil
	},
}

// The Batch Create standard method should have repeated standard create request
// message field.
var requestsField = &lint.MessageRule{
	Name:   lint.NewRuleName("core", "0233", "request-message", "name-field"),
	URI:    "https://aip.dev/233#request-message",
	OnlyIf: isBatchCreateRequestMessage,
	LintMessage: func(m *desc.MessageDescriptor) (problems []lint.Problem) {
		// Rule check: Establish that a "requests" field is present.
		requests := m.FindFieldByName("requests")

		// Rule check: Ensure that the "requests" field is existed.
		if requests == nil {
			return []lint.Problem{{
				Message:    fmt.Sprintf(`Message %q has no "requests" field`, m.GetName()),
				Descriptor: m,
			}}
		}

		// Rule check: Ensure that the standard create request message field "requests" is repeated.
		if !requests.IsRepeated() {
			problems = append(problems, lint.Problem{
				Message:    `The "requests" field should be repeated`,
				Descriptor: requests,
			})
		}

		// Rule check: Ensure that the standard create request message field is the
		// correct type. Note: Use m.GetName()[11:len(m.GetName())-7]) to retrieve
		// the resource name from the the batch create request, for example:
		// "BatchCreateBooksRequest" -> "Books"
		rightTypeName := fmt.Sprintf("Create%sRequest", pluralize.NewClient().Singular(m.GetName()[11:len(m.GetName())-7]))
		if requests.GetMessageType() == nil || requests.GetMessageType().GetName() != rightTypeName {
			problems = append(problems, lint.Problem{
				Message:    fmt.Sprintf(`The "requests" field on Batch Create Request should be a %q type`, rightTypeName),
				Descriptor: requests,
			})
		}
		return
	},
}

// The Batch Create response message should have resource field.
var resourceField = &lint.MessageRule{
	Name:   lint.NewRuleName("core", "0233", "response-message", "resource-field"),
	URI:    "https://aip.dev/233#response-message",
	OnlyIf: isBatchCreateResponseMessage,
	LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
		// the singular form the resource name, the first letter is Capitalized.
		// Note: Use m.GetName()[11 : len(m.GetName())-9] to retrieve the resource
		// name from the the batch create response, for example:
		// "BatchCreateBooksResponse" -> "Books"
		resourceMsgName := pluralize.NewClient().Singular(m.GetName()[11 : len(m.GetName())-9])

		for _, fieldDesc := range m.GetFields() {
			if msgDesc := fieldDesc.GetMessageType(); msgDesc != nil && msgDesc.GetName() == resourceMsgName {
				if !fieldDesc.IsRepeated() {
					return []lint.Problem{{
						Message:    fmt.Sprintf("The %q type field on Batch Create Response message should be repeated", msgDesc.GetName()),
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
