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

package aip0133

import (
	"fmt"
	"strings"

	"bitbucket.org/creachadair/stringset"
	"github.com/googleapis/api-linter/v2/lint"
	"github.com/googleapis/api-linter/v2/rules/internal/utils"
	"github.com/stoewer/go-strcase"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// The create request message should not have unrecognized fields.
var requestRequiredFields = &lint.MethodRule{
	Name:   lint.NewRuleName(133, "request-required-fields"),
	OnlyIf: utils.IsCreateMethodWithResolvedReturnType,
	LintMethod: func(m protoreflect.MethodDescriptor) []lint.Problem {
		ot := utils.GetResponseType(m)
		if ot == nil {
			return nil
		}
		r := utils.GetResource(ot)
		resourceMsgName := utils.GetResourceSingular(r)

		// Rule check: Establish that there are no unexpected fields.
		allowedRequiredFields := stringset.New(
			"parent",
			fmt.Sprintf("%s_id", strings.ToLower(strcase.SnakeCase(resourceMsgName))),
		)

		problems := []lint.Problem{}
		for i := 0; i < m.Input().Fields().Len(); i++ {
			f := m.Input().Fields().Get(i)
			if !utils.GetFieldBehavior(f).Contains("REQUIRED") {
				continue
			}
			// Skip the check with the field that is the resource, which for
			// Standard Create, is the output type.
			if t := f.Message(); t != nil && t.Name() == ot.Name() {
				continue
			}
			// Iterate remaining fields. If they're not in the allowed list,
			// add a problem.
			if !allowedRequiredFields.Contains(string(f.Name())) {
				problems = append(problems, lint.Problem{
					Message:    fmt.Sprintf("Create RPCs must only require fields explicitly described in AIPs, not %q.", f.Name()),
					Descriptor: f,
				})
			}
		}
		return problems
	},
}
