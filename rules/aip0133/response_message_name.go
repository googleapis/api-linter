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

package aip0133

import (
	"fmt"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Create method should use the resource as the output message
var outputName = &lint.MethodRule{
	Name:   lint.NewRuleName(133, "response-message-name"),
	OnlyIf: utils.IsCreateMethod,
	LintMethod: func(m protoreflect.MethodDescriptor) []lint.Problem {
		want := utils.GetResourceMessageName(m, "Create")

		// Load the response type, resolving the
		// `google.longrunning.OperationInfo.response_type` if necessary.
		resp := utils.GetResponseType(m)
		if resp == nil {
			// If we can't resolve it, let the AIP-151 rule warn about this.
			return nil
		}
		got := resp.Name()

		// Rule check: Establish that for methods such as `CreateFoo`, the response
		// message should be named `Foo`
		//
		// Note: If `got` is empty string, this is an unannotated LRO.
		// The AIP-151 rule will whine about that, and this rule should not as it
		// would be confusing.
		if got != want && got != "" {
			return []lint.Problem{{
				Message: fmt.Sprintf(
					"Create RPCs should have the corresponding resource as the response message, such as %q.",
					want,
				),
				Suggestion: want,
				Descriptor: m,
				Location:   locations.MethodResponseType(m),
			}}
		}

		return nil
	},
}
