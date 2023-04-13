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

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

// Batch Update method should have a properly named Response message.
var responseMessageName = &lint.MethodRule{
	Name:   lint.NewRuleName(234, "response-message-name"),
	OnlyIf: isBatchUpdateMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		// Proper response name should be the concatenation of the method name and
		// "Response"
		want := m.GetName() + "Response"

		// If this is an LRO, then use the annotated response type instead of
		// the actual RPC return type.
		got := m.GetOutputType().GetName()
		if utils.IsOperation(m.GetOutputType()) {
			got = utils.GetOperationInfo(m).GetResponseType()
		}

		// Rule check: Establish that for methods such as `BatchUpdateFoos`, the
		// response message should be named `BatchUpdateFoosResponse`
		//
		// Note: If `got` is empty string, this is an unannotated LRO.
		// The AIP-151 rule will whine about that, and this rule should not as it
		// would be confusing.
		if got != want && got != "" {
			return []lint.Problem{{
				Message: fmt.Sprintf(
					"Batch Update RPCs should have a properly named response message %q, but not %q",
					want, got,
				),
				Descriptor: m,
				Location:   locations.MethodResponseType(m),
				Suggestion: want,
			}}
		}

		return nil
	},
}
