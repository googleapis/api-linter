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

package aip0134

import (
	"fmt"

	"bitbucket.org/creachadair/stringset"
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

// The update request message should not have unrecognized fields.
var requestRequiredFields = &lint.MessageRule{
	Name:   lint.NewRuleName(134, "request-required-fields"),
	OnlyIf: isUpdateRequestMessage,
	LintMessage: func(m *desc.MessageDescriptor) (problems []lint.Problem) {
		resourceMsgName := extractResource(m.GetName())

		allowedRequiredFields := stringset.New("update_mask")

		for _, f := range m.GetFields() {
			if !utils.GetFieldBehavior(f).Contains("REQUIRED") {
				continue
			}
			// Skip the check with the field that is the body.
			if t := f.GetMessageType(); t != nil && t.GetName() == resourceMsgName {
				continue
			}
			// add a problem.
			if !allowedRequiredFields.Contains(string(f.GetName())) {
				problems = append(problems, lint.Problem{
					Message:    fmt.Sprintf("Update RPCs must only require fields explicitly described in AIPs, not %q.", f.GetName()),
					Descriptor: f,
				})
			}
		}

		return problems
	},
}
