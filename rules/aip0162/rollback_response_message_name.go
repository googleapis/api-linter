// Copyright 2021 Google LLC
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

package aip0162

import (
	"fmt"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/jhump/protoreflect/desc"
)

var rollbackResponseMessageName = &lint.MethodRule{
	Name:   lint.NewRuleName(162, "rollback-response-message-name"),
	OnlyIf: isRollbackMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		// Rule check: Establish that for methods such as `RollbackBook`, the response
		// message is `Book`.
		want := rollbackMethodRegexp.FindStringSubmatch(m.GetName())[1]
		got := m.GetOutputType().GetName()

		// Return a problem if we did not get the expected return name.
		if got != want {
			return []lint.Problem{{
				Message: fmt.Sprintf(
					"Rollback RPCs should have response message type %q, not %q.",
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
