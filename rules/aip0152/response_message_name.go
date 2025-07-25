// Copyright 2020 Google LLC
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

package aip0152

import (
	"fmt"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var responseMessageName = &lint.MethodRule{
	Name:   lint.NewRuleName(152, "response-message-name"),
	OnlyIf: isRunMethod,
	LintMethod: func(m protoreflect.MethodDescriptor) []lint.Problem {
		if m.Output().FullName() != "google.longrunning.Operation" {
			return []lint.Problem{{
				Message:    "Run methods should use an LRO.",
				Descriptor: m,
				Location:   locations.MethodResponseType(m),
				Suggestion: "google.longrunning.Operation",
			}}
		}

		// Rule check: Establish that for methods such as `RunWriteBookJob`, the
		// response message should be named `RunWriteBookJobRequest`.
		//
		// Note: If `got` is empty string, this is an unannotated LRO.
		// The AIP-151 rule will whine about that, and this rule should not as it
		// would be confusing.
		got := utils.GetOperationInfo(m).GetResponseType()
		want := string(m.Name()) + "Response"
		if got != want && got != "" {
			return []lint.Problem{{
				Message: fmt.Sprintf(
					"Run RPCs should have response message type of `%s`, not %q.",
					want, got,
				),
				Suggestion: want,
				Descriptor: m,
				Location:   locations.MethodResponseType(m),
			}}
		}
		return nil
	},
}
