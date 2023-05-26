// Copyright 2023 Google LLC
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

	"bitbucket.org/creachadair/stringset"
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

// The delete request message should not have unrecognized fields.
var requestRequiredFields = &lint.MessageRule{
	Name:   lint.NewRuleName(135, "request-required-fields"),
	OnlyIf: isDeleteRequestMessage,
	LintMessage: func(m *desc.MessageDescriptor) (problems []lint.Problem) {
		// Rule check: Establish that there are no unexpected fields.
		allowedRequiredFields := stringset.New("name")

		for _, f := range m.GetFields() {
			if !utils.GetFieldBehavior(f).Contains("REQUIRED") {
				continue
			}
			// Iterate remaining fields. If they're not in the allowed list,
			// add a problem.
			if !allowedRequiredFields.Contains(string(f.GetName())) {
				problems = append(problems, lint.Problem{
					Message:    fmt.Sprintf("Delete RPCs must only require fields explicitly described in AIPs, not %q.", f.GetName()),
					Descriptor: f,
				})
			}
		}

		return problems
	},
}
