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

package aip0135

import (
	"fmt"
	"strings"

	"bitbucket.org/creachadair/stringset"
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

// Delete messages should use google.protobuf.Empty,
// google.longrunning.Operation, or the resource itself as the response
// message.
var responseMessageName = &lint.MethodRule{
	Name:   lint.NewRuleName(135, "response-message-name"),
	OnlyIf: isDeleteMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		// Rule check: Establish that for methods such as `DeleteFoo`, the response
		// message is `google.protobuf.Empty` or `Foo`.
		got := m.GetOutputType().GetName()
		if stringset.New("Empty", "Operation").Contains(got) {
			got = m.GetOutputType().GetFullyQualifiedName()
		}
		want := stringset.New(
			"google.protobuf.Empty",
			strings.Replace(m.GetName(), "Delete", "", 1),
		)

		// If the return type is an Operation, use the annotated response type.
		if got == "google.longrunning.Operation" {
			got = utils.GetOperationInfo(m).GetResponseType()
		}

		// If we did not get a permitted value, return a problem.
		//
		// Note: If `got` is empty string, this is an unannotated LRO.
		// The AIP-151 rule will whine about that, and this rule should not as it
		// would be confusing.
		if !want.Contains(got) && got != "" {
			return []lint.Problem{{
				Message: fmt.Sprintf(
					"Delete RPCs should have response message type of Empty or the resource, not %q.",
					got,
				),
				Suggestion: "google.protobuf.Empty",
				Descriptor: m,
			}}
		}

		return nil
	},
}
