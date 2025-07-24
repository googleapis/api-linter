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

package aip0132

import (
	"fmt"

	"bitbucket.org/creachadair/stringset"
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// The list request message should not have unrecognized fields.
var requestRequiredFields = &lint.MessageRule{
	Name:   lint.NewRuleName(132, "request-required-fields"),
	OnlyIf: utils.IsListRequestMessage,
	LintMessage: func(m protoreflect.MessageDescriptor) (problems []lint.Problem) {
		// Rule check: Establish that there are no unexpected fields.
		allowedRequiredFields := stringset.New("parent")

		for _, f := range m.Fields() {
			if !utils.GetFieldBehavior(f).Contains("REQUIRED") {
				continue
			}
			// add a problem.
			if !allowedRequiredFields.Contains(string(f.Name())) {
				problems = append(problems, lint.Problem{
					Message:    fmt.Sprintf("List RPCs must only require fields explicitly described in AIPs, not %q.", f.Name()),
					Descriptor: f,
				})
			}
		}

		return problems
	},
}
