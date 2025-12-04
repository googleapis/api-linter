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

// The create request message should have standardized field types for required fields.
var requestRequiredFieldsType = &lint.MethodRule{
	Name:   lint.NewRuleName(133, "request-required-fields-type"),
	OnlyIf: utils.IsCreateMethodWithResolvedReturnType,
	LintMethod: func(m protoreflect.MethodDescriptor) []lint.Problem {
		ot := utils.GetResponseType(m)
		var resourceMsgName string

		// Try to get resource from response type
		if r := utils.GetResource(ot); r != nil {
			resourceMsgName = utils.GetResourceSingular(r)
		}

		// If we can't get a resource from the response types, we can try to infer it from the method
		if resourceMsgName == "" {
			if noun := utils.GetResourceMessageName(m, "Create"); noun != "" {
				resourceMsgName = noun
			}
		}

		problems := []lint.Problem{}
		reqFields := m.Input().Fields()

		// We can check for each field explicitly since we can infer what each field should be

		// Check for `parent`.
		if parentField := reqFields.ByName("parent"); parentField != nil {
			if parentField.Kind() != protoreflect.StringKind {
				problems = append(problems, lint.Problem{
					Message:    `The field "parent" must be of type string.`,
					Descriptor: parentField,
				})
			}
		}

		snakeResourceName := strings.ToLower(strcase.SnakeCase(resourceMsgName))
		resourceIdFieldName := protoreflect.Name(fmt.Sprintf("%s_id", snakeResourceName))
		resourceFieldName := protoreflect.Name(snakeResourceName)

		// Check for `<resource>_id`.
		if idField := reqFields.ByName(resourceIdFieldName); idField != nil {
			if idField.Kind() != protoreflect.StringKind {
				problems = append(problems, lint.Problem{
					Message:    fmt.Sprintf("The field %q must be of type string.", idField.Name()),
					Descriptor: idField,
				})
			}
		}

		// Check for `<resource>`.
		if resourceField := reqFields.ByName(resourceFieldName); resourceField != nil {
			if resourceField.Kind() != protoreflect.MessageKind {
				problems = append(problems, lint.Problem{
					Message:    fmt.Sprintf("The field %q must be of type message.", resourceField.Name()),
					Descriptor: resourceField,
				})
			}
		}

		return problems
	},
}
