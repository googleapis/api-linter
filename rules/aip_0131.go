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

package rules

import (
	"fmt"
	"strings"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/descriptor"
	p "google.golang.org/protobuf/reflect/protoreflect"
)

func init() {
	registerRules(
		checkGetMethodNameField(),
		checkRequestMessageName(),
		checkUnknownFields(),
	)
}

// The Get standard method should only have expected fields.
func checkGetMethodNameField() lint.Rule {
	return &descriptor.CallbackRule{
		RuleInfo: lint.RuleInfo{
			Name:         lint.NewRuleName("core", "0131", "name_field"),
			Description:  "check that a name field is present",
			RequestTypes: []lint.RequestType{lint.ProtoRequest},
		},
		Callback: descriptor.Callbacks{
			MethodCallback: func(m p.MethodDescriptor, s lint.DescriptorSource) (problems []lint.Problem, err error) {
				// We only care about Get- methods for the purpose of this rule;
				// ignore everything else.
				if !strings.HasPrefix(string(m.Name()), "Get") {
					return
				}

				// Rule check: Establish that a name field is present.
				nameField := m.Input().Fields().ByName("name")
				if nameField == nil {
					problems = append(problems, lint.Problem{
						Message:    fmt.Sprintf("method %q has no `name` field", m.Name()),
						Descriptor: m.Input(),
					})
					return problems, nil
				}

				// Rule check: Establish that the name field is a string.
				if nameField.Kind() != p.StringKind {
					problems = append(problems, lint.Problem{
						Message:    "`name` field on Get RPCs should be a string",
						Descriptor: nameField,
					})
				}

				return problems, nil
			},
		},
	}
}

func checkUnknownFields() lint.Rule {
	return &descriptor.CallbackRule{
		RuleInfo: lint.RuleInfo{
			Name:         lint.NewRuleName("core", "0131", "unknown_fields"),
			Description:  "check that there are no unknown fields",
			RequestTypes: []lint.RequestType{lint.ProtoRequest},
		},
		Callback: descriptor.Callbacks{
			MethodCallback: func(m p.MethodDescriptor, s lint.DescriptorSource) (problems []lint.Problem, err error) {
				// We only care about Get- methods for the purpose of this rule;
				// ignore everything else.
				if !strings.HasPrefix(string(m.Name()), "Get") {
					return
				}

				// Rule check: Establish that there are no other fields besides `name`.
				fields := m.Input().Fields()
				for i := 0; i < fields.Len(); i++ {
					field := fields.Get(i)
					if field.Name() != "name" {
						problems = append(problems, lint.Problem{
							Message:    "Get RPCs should not have fields other than `name`.",
							Descriptor: field,
						})
					}
				}

				return problems, nil
			},
		},
	}
}

func checkRequestMessageName() lint.Rule {
	return &descriptor.CallbackRule{
		RuleInfo: lint.RuleInfo{
			Name:         lint.NewRuleName("core", "0131", "request_message_name"),
			Description:  "check that Get RPCs have appropriate request messages",
			RequestTypes: []lint.RequestType{lint.ProtoRequest},
		},
		Callback: descriptor.Callbacks{
			MethodCallback: func(m p.MethodDescriptor, s lint.DescriptorSource) (problems []lint.Problem, err error) {
				// We only care about Get- methods for the purpose of this rule;
				// ignore everything else.
				methodName := string(m.Name())
				if !strings.HasPrefix(methodName, "Get") {
					return
				}

				// Rule check: Establish that for methods such as `GetFoo`, the request
				// message is named `GetFooRequest`.
				requestMessageName := string(m.Input().Name())
				if requestMessageName != methodName+"Request" {
					problems = append(problems, lint.Problem{
						Message:    fmt.Sprintf("Get RPCs should have a request message named after the RPC, such as %q.", methodName+"Request"),
						Suggestion: methodName + "Request",
						Descriptor: m,
					})
				}

				return problems, nil
			},
		},
	}
}
