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

	"github.com/commure/api-linter/lint"
	"github.com/commure/api-linter/locations"
	"github.com/gertd/go-pluralize"
	"github.com/jhump/protoreflect/desc"
)

// The Batch Create standard method should have repeated standard create request
// message field.
var requestRequestsField = &lint.MessageRule{
	Name:   lint.NewRuleName(233, "request-requests-field"),
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
		// correct type. Note: Retrieve the resource name from the the batch create
		// request, for example: "BatchCreateBooksRequest" -> "Books"
		rightTypeName := fmt.Sprintf("Create%sRequest",
			pluralize.NewClient().Singular(strings.TrimPrefix(strings.TrimSuffix(m.GetName(), "Request"), "BatchCreate")))
		if requests.GetMessageType() == nil || requests.GetMessageType().GetName() != rightTypeName {
			problems = append(problems, lint.Problem{
				Message:    fmt.Sprintf(`The "requests" field on Batch Create Request should be a %q type`, rightTypeName),
				Descriptor: requests,
				Location:   locations.FieldType(requests),
				Suggestion: rightTypeName,
			})
		}
		return
	},
}
