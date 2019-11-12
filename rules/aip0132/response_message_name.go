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

package aip0132

import (
	"fmt"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/jhump/protoreflect/desc"
)

// List messages should use a `ListFoosResponse` response message.
var responseMessageName = &lint.MethodRule{
	Name:   lint.NewRuleName(132, "response-message-name"),
	OnlyIf: isListMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		// Rule check: Establish that for methods such as `ListFoos`, the response
		// message is named `ListFoosResponse`.
		if got, want := m.GetOutputType().GetName(), m.GetName()+"Response"; got != want {
			return []lint.Problem{{
				Message: fmt.Sprintf(
					"List RPCs should have a response message named after the RPC, such as %q.",
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
