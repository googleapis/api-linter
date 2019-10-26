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

	"github.com/googleapis/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"
	"github.com/kjk/inflect"
)

// The Batch Get request message should have parent field.
var parentField = &lint.MessageRule{
	Name:   lint.NewRuleName("core", "0231", "request-message", "parent-field"),
	URI:    "https://aip.dev/231#request-message",
	OnlyIf: isBatchGetRequestMessage,
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
				Message:    "`parent` field on create request message must be a string",
				Descriptor: parentField,
			}}
		}

		return nil
	},
}

// The Batch Get standard method should have repeated name field or repeated
// standard get request message field, but the latter one is not suggested.
var namesField = &lint.MessageRule{
	Name:   lint.NewRuleName("core", "0231", "request-message", "name-field"),
	URI:    "https://aip.dev/231#request-message",
	OnlyIf: isBatchGetRequestMessage,
	LintMessage: func(m *desc.MessageDescriptor) (problems []lint.Problem) {
		// Rule check: Establish that a name field is present.
		names := m.FindFieldByName("names")
		getReqMsg := m.FindFieldByName("requests")

		// Rule check: Ensure that the names field is existed.
		if names == nil && getReqMsg == nil {
			problems = append(problems, lint.Problem{
				Message:    fmt.Sprintf(`Message %q has no "names" field`, m.GetName()),
				Descriptor: m,
			})
		}

		// Rule check: Ensure that only the suggested names field is existed.
		if names != nil && getReqMsg != nil {
			problems = append(problems, lint.Problem{
				Message:    fmt.Sprintf(`Message %q should delete "requests" field, only keep the "names" field`, m.GetName()),
				Descriptor: getReqMsg,
			})
		}

		// Rule check: Ensure that the names field is repeated.
		if names != nil && !names.IsRepeated() {
			problems = append(problems, lint.Problem{
				Message:    `The "names"" field should be repeated`,
				Descriptor: names,
			})
		}

		// Rule check: Ensure that the names field is the correct type.
		if names != nil && names.GetType() != builder.FieldTypeString().GetType() {
			problems = append(problems, lint.Problem{
				Message:    `"names" field on Batch Get Request should be a "string" type`,
				Descriptor: names,
			})
		}

		// Rule check: Ensure that the standard get request message field is repeated.
		if getReqMsg != nil && !getReqMsg.IsRepeated() {
			problems = append(problems, lint.Problem{
				Message:    `The "requests" field should be repeated`,
				Descriptor: getReqMsg,
			})
		}
		// Rule check: Ensure that the standard get request message field is the correct type.
		rightTypeName := fmt.Sprintf("Get%sRequest", inflect.ToSingular(m.GetName()[8:len(m.GetName())-7]))
		if getReqMsg != nil && (getReqMsg.GetMessageType() == nil || getReqMsg.GetMessageType().GetName() != rightTypeName) {
			problems = append(problems, lint.Problem{
				Message:    fmt.Sprintf(`The "requests" field on Batch Get Request should be a %q type`, rightTypeName),
				Descriptor: getReqMsg,
			})
		}
		return
	},
}

// The Batch Get response message should have resource field.
var resourceField = &lint.MessageRule{
	Name:   lint.NewRuleName("core", "0231", "response-message", "resource-field"),
	URI:    "https://aip.dev/231#response-message",
	OnlyIf: isBatchGetResponseMessage,
	LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
		// the singular form the resource name, the first letter is Capitalized
		resourceMsgName := inflect.ToSingular(m.GetName()[8 : len(m.GetName())-9])

		for _, fieldDesc := range m.GetFields() {
			if msgDesc := fieldDesc.GetMessageType(); msgDesc != nil && msgDesc.GetName() == resourceMsgName {
				if !fieldDesc.IsRepeated() {
					return []lint.Problem{{
						Message:    fmt.Sprintf("The %q type field on Batch Get Response message should be repeated", msgDesc.GetName()),
						Descriptor: fieldDesc,
					}}
				} else {
					return nil
				}
			}
		}

		// Rule check: Establish that a resource field must be included.
		return []lint.Problem{{
			Message:    fmt.Sprintf("Message %q has no %q type field", m.GetName(), resourceMsgName),
			Descriptor: m,
		}}
	},
}
