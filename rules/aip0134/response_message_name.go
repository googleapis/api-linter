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

package aip0134

import (
	"fmt"
	"strings"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

// Update methods should use the resource as the response message
var responseMessageName = &lint.MethodRule{
	Name:   lint.NewRuleName(134, "response-message-name"),
	OnlyIf: utils.IsUpdateMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		// Rule check: Establish that for methods such as `UpdateFoo`, the response
		// message is `Foo` or `google.longrunning.Operation`.
		want := strings.Replace(m.GetName(), "Update", "", 1)

		// Load the response type, resolving the
		// `google.longrunning.OperationInfo.response_type` if necessary.
		resp := utils.GetResponseType(m)
		if resp == nil {
			// If we can't resolve it, let the AIP-151 rule warn about this.
			return nil
		}
		got := resp.GetName()

		// Return a problem if we did not get the expected return name.
		//
		// Note: If `got` is empty string, this is an unannotated LRO.
		// The AIP-151 rule will whine about that, and this rule should not as it
		// would be confusing.
		if got != want && got != "" {
			return []lint.Problem{{
				Message: fmt.Sprintf(
					"Update RPCs should have response message type %q, not %q.",
					want,
					got,
				),
				Suggestion: want,
				Descriptor: m,
				Location:   locations.MethodResponseType(m),
			}}
		}
		return nil
	},
}
