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

		var resourceMsgName string
		if r := utils.GetResource(ot); r != nil {
			resourceMsgName = utils.GetResourceSingular(r)
		}

		if resourceMsgName == "" {
			if noun := utils.GetResourceMessageName(m, "Create"); noun != "" {
				resourceMsgName = noun
			}
		}

		snakeResourceName := strings.ToLower(strcase.SnakeCase(resourceMsgName))

		// Rule check: Establish that there are no unexpected fields.
		allowedRequiredFields := map[string]protoreflect.Kind{
			"parent":                                protoreflect.StringKind,
			fmt.Sprintf("%s_id", snakeResourceName): protoreflect.StringKind,
			snakeResourceName:                       protoreflect.MessageKind,
		}

		problems := []lint.Problem{}
		for i := 0; i < m.Input().Fields().Len(); i++ {
			f := m.Input().Fields().Get(i)
			if !utils.GetFieldBehavior(f).Contains("REQUIRED") {
				continue
			}

			// Iterate remaining fields. If they're not in the allowed list,
			// add a problem.
			expectedKind, allowed := allowedRequiredFields[string(f.Name())]
			if !allowed {
				problems = append(problems, lint.Problem{
					Message:    fmt.Sprintf("Create RPCs must only require fields explicitly described in AIPs, not %q.", f.Name()),
					Descriptor: f,
				})
			} else if f.Kind() != expectedKind {
				problems = append(problems, lint.Problem{
					Message:    fmt.Sprintf("The required field %q must be of type %v.", f.Name(), expectedKind),
					Descriptor: f,
				})
			}
		}
		return problems
	},
}
