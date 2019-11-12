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
	"github.com/googleapis/api-linter/locations"
	"github.com/jhump/protoreflect/desc"
)

// Batch Get method should have a properly named Request message.
var inputName = &lint.MethodRule{
	Name:   lint.NewRuleName(231, "request-message-name"),
	OnlyIf: isBatchGetMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		// Rule check: Establish that for methods such as `BatchGetFoos`, the request
		// message is named `BatchGetFoosRequest`.
		if got, want := m.GetInputType().GetName(), m.GetName()+"Request"; got != want {
			return []lint.Problem{{
				Message: fmt.Sprintf(
					"Batch Get RPCs should have a properly named request message %q, but not %q",
					want, got,
				),
				Suggestion: want,
				Descriptor: m,
				Location:   locations.MethodRequestType(m),
			}}
		}

		return nil
	},
}
