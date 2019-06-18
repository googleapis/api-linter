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
		checkListRequestMessageName(),
		checkListRequestMessageParentField(),
		checkListRequestMessageUnknownFields(),
		checkListResponseMessageName(),
	)
}

// Get messages should have a properly named Request message.
func checkListRequestMessageName() lint.Rule {
	return &descriptor.CallbackRule{
		RuleInfo: lint.RuleInfo{
			Name:         lint.NewRuleName("core", "0132", "request_message", "name"),
			Description:  "check that List RPCs have appropriate request messages",
			RequestTypes: []lint.RequestType{lint.ProtoRequest},
			URI:          "https://aip.dev/132#guidance",
		},
		Callback: descriptor.Callbacks{
			MethodCallback: func(m p.MethodDescriptor, s lint.DescriptorSource) (problems []lint.Problem, err error) {
				// We only care about List- methods for the purpose of this rule;
				// ignore everything else.
				if !isListMethod(m) {
					return
				}

				// Rule check: Establish that for methods such as `GetFoo`, the request
				// message is named `GetFooRequest`.
				methodName := string(m.Name())
				requestMessageName := string(m.Input().Name())
				if requestMessageName != methodName+"Request" {
					problems = append(problems, lint.Problem{
						Message: fmt.Sprintf(
							"List RPCs should have a request message named after the RPC, such as %q.",
							methodName+"Request",
						),
						Suggestion: methodName + "Request",
						Descriptor: m,
					})
				}

				return problems, nil
			},
		},
	}
}

// The List standard method should contain a parent field.
func checkListRequestMessageParentField() lint.Rule {
	return &descriptor.CallbackRule{
		RuleInfo: lint.RuleInfo{
			Name:         lint.NewRuleName("core", "0132", "request_message", "parent_field"),
			Description:  "check that a parent field is present",
			RequestTypes: []lint.RequestType{lint.ProtoRequest},
			URI:          "https://aip.dev/132#request-message",
		},
		Callback: descriptor.Callbacks{
			MethodCallback: func(m p.MethodDescriptor, s lint.DescriptorSource) (problems []lint.Problem, err error) {
				// We only care about List- methods for the purpose of this rule;
				// ignore everything else.
				if !isListMethod(m) {
					return
				}

				// Rule check: Establish that a name field is present.
				parentField := m.Input().Fields().ByName("parent")
				if parentField == nil {
					problems = append(problems, lint.Problem{
						Message:    fmt.Sprintf("method %q has no `parent` field", m.Name()),
						Descriptor: m.Input(),
					})
					return problems, nil
				}

				// Rule check: Establish that the name field is a string.
				if parentField.Kind() != p.StringKind {
					problems = append(problems, lint.Problem{
						Message:    "`parent` field on List RPCs should be a string",
						Descriptor: parentField,
					})
				}

				return problems, nil
			},
		},
	}
}

// Get methods should not have unrecognized fields.
func checkListRequestMessageUnknownFields() lint.Rule {
	return &descriptor.CallbackRule{
		RuleInfo: lint.RuleInfo{
			Name:         lint.NewRuleName("core", "0132", "request_message", "unknown_fields"),
			Description:  "check that there are no unknown fields",
			RequestTypes: []lint.RequestType{lint.ProtoRequest},
			URI:          "https://aip.dev/131#request-message",
		},
		Callback: descriptor.Callbacks{
			MethodCallback: func(m p.MethodDescriptor, s lint.DescriptorSource) (problems []lint.Problem, err error) {
				// We only care about List- methods for the purpose of this rule;
				// ignore everything else.
				if !isListMethod(m) {
					return
				}

				// Rule check: Establish that there are no unexpected fields.
				allowedFields := map[string]struct{}{
					"parent":       struct{}{}, // AIP-131
					"page_size":    struct{}{}, // AIP-158
					"page_token":   struct{}{}, // AIP-158
					"filter":       struct{}{}, // AIP-132
					"order_by":     struct{}{}, // AIP-132
					"group_by":     struct{}{}, // Nowhere yet, but permitted
					"show_deleted": struct{}{}, // AIP-135
					"read_mask":    struct{}{}, // AIP-157
					"view":         struct{}{}, // AIP-157
				}
				fields := m.Input().Fields()
				for i := 0; i < fields.Len(); i++ {
					field := fields.Get(i)
					if _, ok := allowedFields[string(field.Name())]; !ok {
						problems = append(problems, lint.Problem{
							Message:    "List RPCs should only contain fields explicitly described in AIPs.",
							Descriptor: field,
						})
					}
				}

				return problems, nil
			},
		},
	}
}

// Get messages should use the resource as the respose message
func checkListResponseMessageName() lint.Rule {
	return &descriptor.CallbackRule{
		RuleInfo: lint.RuleInfo{
			Name:         lint.NewRuleName("core", "0132", "response_message", "name"),
			Description:  "check that List RPCs have appropriate response messages",
			RequestTypes: []lint.RequestType{lint.ProtoRequest},
			URI:          "https://aip.dev/131#guidance",
		},
		Callback: descriptor.Callbacks{
			MethodCallback: func(m p.MethodDescriptor, s lint.DescriptorSource) (problems []lint.Problem, err error) {
				// We only care about List- methods for the purpose of this rule;
				// ignore everything else.
				if !isListMethod(m) {
					return
				}

				// Rule check: Establish that for methods such as `ListFoos`, the response
				// message is named `ListFoosResponse`.
				methodName := string(m.Name())
				responseMessageName := string(m.Output().Name())
				if responseMessageName != methodName+"Response" {
					problems = append(problems, lint.Problem{
						Message: fmt.Sprintf(
							"List RPCs should have a response message named after the RPC, such as %q.",
							methodName+"Response",
						),
						Suggestion: methodName + "Response",
						Descriptor: m,
					})
				}

				return problems, nil
			},
		},
	}
}

// Return true if this is a List method, false otherwise.
func isListMethod(m p.MethodDescriptor) bool {
	return strings.HasPrefix(string(m.Name()), "List")
}
