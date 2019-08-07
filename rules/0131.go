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
	"regexp"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/descriptor"
	"github.com/golang/protobuf/proto"
	"google.golang.org/genproto/googleapis/api/annotations"
	desc "google.golang.org/protobuf/types/descriptorpb"
	p "google.golang.org/protobuf/reflect/protoreflect"
)

func init() {
	registerRules(
		checkGetRequestMessageName(),
		checkGetRequestMessageNameField(),
		checkGetRequestMessageUnknownFields(),
		checkGetResponseMessageName(),
	)
}

// Get messages should have a properly named Request message.
func checkGetRequestMessageName() lint.Rule {
	return &descriptor.CallbackRule{
		RuleInfo: lint.RuleInfo{
			Name:         lint.NewRuleName("core", "0131", "request_message", "name"),
			Description:  "check that Get RPCs have appropriate request messages",
			RequestTypes: []lint.RequestType{lint.ProtoRequest},
			URI:          "https://aip.dev/131#guidance",
		},
		Callback: descriptor.Callbacks{
			MethodCallback: func(m p.MethodDescriptor, s lint.DescriptorSource) (problems []lint.Problem, err error) {
				// We only care about Get methods for the purpose of this rule;
				// ignore everything else.
				if !isGetMethod(m) {
					return
				}

				// Rule check: Establish that for methods such as `GetFoo`, the request
				// message is named `GetFooRequest`.
				methodName := string(m.Name())
				requestMessageName := string(m.Input().Name())
				correctRequestMessageName := methodName + "Request"
				if requestMessageName != correctRequestMessageName {
					problems = append(problems, lint.Problem{
						Message: fmt.Sprintf(
							"Get RPCs should have a request message named after the RPC, such as %q.",
							correctRequestMessageName,
						),
						Suggestion: correctRequestMessageName,
						Descriptor: m,
					})
				}

				return problems, nil
			},
		},
	}
}

// Get messages should have a properly named Request message.
func checkGetURI() lint.Rule {
	return &descriptor.CallbackRule{
		RuleInfo: lint.RuleInfo{
			Name:         lint.NewRuleName("core", "0131", "uri"),
			Description:  "check that Get RPCs have an appropriate URI",
			RequestTypes: []lint.RequestType{lint.ProtoRequest},
			URI:          "https://aip.dev/131#guidance",
		},
		Callback: descriptor.Callbacks{
			MethodCallback: func(m p.MethodDescriptor, s lint.DescriptorSource) (problems []lint.Problem, err error) {
				// We only care about Get methods for the purpose of this rule;
				// ignore everything else.
				if !isGetMethod(m) {
					return
				}

				// Rule check: Establish that for methods such as `GetFoo`, the request
				// URI uses HTTP GET.
				opts := m.Options().(*desc.MethodOptions)
				fmt.Printf("%v\n", opts)
				httpRule, err := proto.GetExtension(
					m.Options().(*desc.MethodOptions),
					annotations.E_Http,
				)
				fmt.Printf("%v\n", httpRule)
				methodName := string(m.Name())
				fmt.Printf("%v\n", proto.WireStartGroup)
				fmt.Printf("%v\n", annotations.E_Http)
				requestMessageName := string(m.Input().Name())
				correctRequestMessageName := methodName + "Request"
				if requestMessageName != correctRequestMessageName {
					problems = append(problems, lint.Problem{
						Message: fmt.Sprintf(
							"Get RPCs should have a request message named after the RPC, such as %q.",
							correctRequestMessageName,
						),
						Suggestion: correctRequestMessageName,
						Descriptor: m,
					})
				}

				return problems, nil
			},
		},
	}
}

// The Get standard method should only have expected fields.
func checkGetRequestMessageNameField() lint.Rule {
	return &descriptor.CallbackRule{
		RuleInfo: lint.RuleInfo{
			Name:         lint.NewRuleName("core", "0131", "request_message", "name_field"),
			Description:  "check that a name field is present",
			RequestTypes: []lint.RequestType{lint.ProtoRequest},
			URI:          "https://aip.dev/131#request-message",
		},
		Callback: descriptor.Callbacks{
			MethodCallback: func(m p.MethodDescriptor, s lint.DescriptorSource) (problems []lint.Problem, err error) {
				// We only care about Get methods for the purpose of this rule;
				// ignore everything else.
				if !isGetMethod(m) {
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

// Get methods should not have unrecognized fields.
func checkGetRequestMessageUnknownFields() lint.Rule {
	return &descriptor.CallbackRule{
		RuleInfo: lint.RuleInfo{
			Name:         lint.NewRuleName("core", "0131", "request_message", "unknown_fields"),
			Description:  "check that there are no unknown fields",
			RequestTypes: []lint.RequestType{lint.ProtoRequest},
			URI:          "https://aip.dev/131#request-message",
		},
		Callback: descriptor.Callbacks{
			MethodCallback: func(m p.MethodDescriptor, s lint.DescriptorSource) (problems []lint.Problem, err error) {
				// We only care about Get methods for the purpose of this rule;
				// ignore everything else.
				if !isGetMethod(m) {
					return
				}

				// Rule check: Establish that there are no unexpected fields.
				allowedFields := map[string]struct{}{
					"name":      {}, // AIP-131
					"read_mask": {}, // AIP-157
					"view":      {}, // AIP-157
				}
				fields := m.Input().Fields()
				for i := 0; i < fields.Len(); i++ {
					field := fields.Get(i)
					if _, ok := allowedFields[string(field.Name())]; !ok {
						problems = append(problems, lint.Problem{
							Message: fmt.Sprintf(
								"Get RPCs should only only contain fields explicitly described in AIPs, not %q.",
								string(field.Name()),
							),
							Descriptor: field,
						})
					}
				}

				return problems, nil
			},
		},
	}
}

// Get messages should use the resource as the response message
func checkGetResponseMessageName() lint.Rule {
	return &descriptor.CallbackRule{
		RuleInfo: lint.RuleInfo{
			Name:         lint.NewRuleName("core", "0131", "response_message", "name"),
			Description:  "check that Get RPCs have appropriate response messages",
			RequestTypes: []lint.RequestType{lint.ProtoRequest},
			URI:          "https://aip.dev/131#guidance",
		},
		Callback: descriptor.Callbacks{
			MethodCallback: func(m p.MethodDescriptor, s lint.DescriptorSource) (problems []lint.Problem, err error) {
				// We only care about Get methods for the purpose of this rule;
				// ignore everything else.
				if !isGetMethod(m) {
					return
				}

				// Rule check: Establish that for methods such as `GetFoo`, the response
				// message is named `Foo`.
				responseMessageName := string(m.Output().Name())
				methodName := string(m.Name())
				if correctResponseMessageName := methodName[3:]; correctResponseMessageName != responseMessageName {
					problems = append(problems, lint.Problem{
						Message: fmt.Sprintf(
							"Get RPCs should have the corresponding resource as the response message, such as %q.",
							correctResponseMessageName,
						),
						Suggestion: correctResponseMessageName,
						Descriptor: m,
					})
				}

				return problems, nil
			},
		},
	}
}

// Return true if this is a AIP-131 Get method, false otherwise.
func isGetMethod(m p.MethodDescriptor) bool {
	methodName := string(m.Name())
	if methodName == "GetIamPolicy" {
		return false
	}
	return regexp.MustCompile("^Get(?:[A-Z]|$)").MatchString(string(m.Name()))
}
