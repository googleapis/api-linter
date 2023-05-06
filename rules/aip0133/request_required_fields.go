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

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"
	"github.com/stoewer/go-strcase"
)

// The create request message should not have unrecognized fields.
var requestRequiredFields = &lint.MessageRule{
	Name:   lint.NewRuleName(133, "request-required-fields"),
	OnlyIf: isCreateRequestMessage,
	LintMessage: func(m *desc.MessageDescriptor) (problems []lint.Problem) {
		resourceMsgName := getResourceMsgNameFromReq(m)

		// Rule check: Establish that there are no unexpected fields.
		allowedRequiredFields := map[string]*builder.FieldType{
			"parent": nil, // AIP-133
			fmt.Sprintf("%s_id", strings.ToLower(strcase.SnakeCase(resourceMsgName))): nil,
		}

		for _, fieldDesc := range m.GetFields() {
			if msgDesc := fieldDesc.GetMessageType(); msgDesc != nil && msgDesc.GetName() == resourceMsgName {
				allowedRequiredFields[fieldDesc.GetName()] = nil // AIP-133
			}
		}

		for _, f := range m.GetFields() {
			if !utils.GetFieldBehavior(f).Contains("REQUIRED") {
				continue
			}
			if _, ok := allowedRequiredFields[string(f.GetName())]; !ok {
				problems = append(problems, lint.Problem{
					Message: fmt.Sprintf(
						"Create RPCs must only require fields explicitly described in AIPs, not %q.",
						f.GetName(),
					),
					Descriptor: f,
				})
			}
		}

		return problems
	},
}
