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
	"github.com/googleapis/api-linter/locations"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// The Batch Update standard method should have repeated standard Update request
// message field.
var requestRequestsField = &lint.MessageRule{
	Name:   lint.NewRuleName(234, "request-requests-field"),
	OnlyIf: isBatchUpdateRequestMessage,
	LintMessage: func(m protoreflect.MessageDescriptor) (problems []lint.Problem) {
		// Rule check: Establish that a "requests" field is present.
		requests := m.FindFieldByName("requests")

		// Rule check: Ensure that the "requests" field is existed.
		if requests == nil {
			return []lint.Problem{{
				Message:    fmt.Sprintf(`Message %q has no "requests" field`, m.Name()),
				Descriptor: m,
			}}
		}

		// Rule check: Ensure that the standard Update request message field "requests" is repeated.
		if !requests.IsRepeated() {
			problems = append(problems, lint.Problem{
				Message:    `The "requests" field should be repeated`,
				Descriptor: requests,
			})
		}

		// Rule check: Ensure that the standard update request message field is the
		// correct type. Note: Retrieve the resource name from the the batch update
		// request, for example: "BatchUpdateBooksRequest" -> "Books"
		rightTypeName := fmt.Sprintf("Update%sRequest",
			pluralize.NewClient().Singular(strings.TrimPrefix(strings.TrimSuffix(m.Name(), "Request"), "BatchUpdate")))
		if requests.GetMessageType() == nil || requests.GetMessageType().Name() != rightTypeName {
			problems = append(problems, lint.Problem{
				Message:    fmt.Sprintf(`The "requests" field on Batch Update Request should be a %q type`, rightTypeName),
				Descriptor: requests,
				Location:   locations.FieldType(requests),
				Suggestion: rightTypeName,
			})
		}
		return
	},
}
