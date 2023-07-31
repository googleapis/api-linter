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

package aip0133

import (
	"fmt"
	"strings"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

// Create methods should have an HTTP body, and the body value should be resource.
var httpBody = &lint.MethodRule{
	Name:   lint.NewRuleName(133, "http-body"),
	OnlyIf: utils.IsCreateMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		resourceMsgName := getResourceMsgName(m)
		resourceFieldName := strings.ToLower(resourceMsgName)
		for _, fieldDesc := range m.GetInputType().GetFields() {
			// when msgDesc is nil, the resource field in the request message is
			// missing. A lint warning for the rule `resourceField` will be generated.
			// For here, we will use the lower case resource message name as default
			if msgDesc := fieldDesc.GetMessageType(); msgDesc != nil && msgDesc.GetName() == resourceMsgName {
				resourceFieldName = fieldDesc.GetName()
			}
		}

		// Establish that HTTP body the RPC should map the resource field name in the request message.
		for _, httpRule := range utils.GetHTTPRules(m) {
			if httpRule.Body == "" {
				// Establish that the RPC should have HTTP body
				return []lint.Problem{{
					Message:    "Post methods should have an HTTP body.",
					Descriptor: m,
					Location:   locations.MethodHTTPRule(m),
				}}
				// When resource field is not set in the request message, the problem
				// will not be triggered by the rule"core::0133::http-body". It will be
				// triggered by another rule
				// "core::0133::request-message::resource-field"
			} else if resourceFieldName != "" && httpRule.Body != resourceFieldName {
				return []lint.Problem{{
					Message: fmt.Sprintf(
						"The content of body %q must map to the resource field %q in the request message",
						httpRule.Body,
						resourceFieldName,
					),
					Descriptor: m,
					Location:   locations.MethodHTTPRule(m),
				}}
			}
		}
		return nil
	},
}
