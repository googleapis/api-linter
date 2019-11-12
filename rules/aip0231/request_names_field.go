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
	"github.com/googleapis/api-linter/locations"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"
)

// The Batch Get standard method should have repeated name field or repeated
// standard get request message field, but the latter one is not suggested.
var namesField = &lint.MessageRule{
	Name:   lint.NewRuleName(231, "request-names-field"),
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
				Message:    `The "names" field should be repeated`,
				Descriptor: names,
				Location:   locations.FieldType(names),
			})
		}

		// Rule check: Ensure that the names field is the correct type.
		if names != nil && names.GetType() != builder.FieldTypeString().GetType() {
			problems = append(problems, lint.Problem{
				Message:    `"names" field on Batch Get Request should be a "string" type`,
				Suggestion: "string",
				Descriptor: names,
				Location:   locations.FieldType(names),
			})
		}

		// Rule check: Ensure that the standard get request message field is repeated.
		if getReqMsg != nil && !getReqMsg.IsRepeated() {
			problems = append(problems, lint.Problem{
				Message:    `The "requests" field should be repeated`,
				Descriptor: getReqMsg,
			})
		}

		// Rule check: Ensure that the standard get request message field is the
		// correct type. Note: Use m.GetName()[8:len(m.GetName())-7]) to retrieve
		// the resource name from the the batch get request, for example:
		// "BatchGetBooksRequest" -> "Books"
		rightTypeName := fmt.Sprintf("Get%sRequest", pluralize.NewClient().Singular(m.GetName()[8:len(m.GetName())-7]))
		if getReqMsg != nil && (getReqMsg.GetMessageType() == nil || getReqMsg.GetMessageType().GetName() != rightTypeName) {
			problems = append(problems, lint.Problem{
				Message:    fmt.Sprintf(`The "requests" field on Batch Get Request should be a %q type`, rightTypeName),
				Descriptor: getReqMsg,
			})
		}
		return
	},
}
