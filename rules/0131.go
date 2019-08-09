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
	p "google.golang.org/protobuf/reflect/protoreflect"
)

func init() {
	registerRules(
		checkGetRequestMessageName(),
		checkGetRequestMessageNameField(),
		checkGetRequestMessageUnknownFields(),
		checkGetResponseMessageName(),
		checkGetURI(),
		checkGetBody(),
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

// Get messages should have an appropriate (google.api.http).get URI.
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

				// Iterate over the HTTP rules.
				for _, httpRule := range getHTTPRules(m) {
					// Rule check: Ensure that the GET HTTP verb is used.
					getURI := httpRule.GetGet()
					if getURI == "" {
						problems = append(problems, lint.Problem{
							Message:    "Get RPCs should use the GET HTTP verb.",
							Descriptor: m,
						})
						return problems, nil
					}

					// Rule check: Ensure that the name variable is included in the URI,
					// and encompasses the entire resource name.
					if !regexp.MustCompile("\\{name=[a-zA-Z/*]+\\}$").MatchString(getURI) {
						problems = append(problems, lint.Problem{
							Message:    "Get RPCs should include the `name` field in the URI. Example: /v1/{name=publishers/*/books/*}",
							Descriptor: m,
						})
					}
				}

				return problems, nil
			},
		},
	}
}

// Get messages should not have an HTTP body.
func checkGetBody() lint.Rule {
	return &descriptor.CallbackRule{
		RuleInfo: lint.RuleInfo{
			Name:         lint.NewRuleName("core", "0131", "body"),
			Description:  "check that Get RPCs have no HTTP body",
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

				// Rule check: Ensure that there is no HTTP body.
				for _, httpRule := range getHTTPRules(m) {
					if httpRule.GetBody() != "" {
						problems = append(problems, lint.Problem{
							Message:    "Get RPCs should not have an HTTP body. Ensure the `body` key in the google.api.http annotation is absent.",
							Descriptor: m,
						})
					}
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
			MessageCallback: func(m p.MessageDescriptor, s lint.DescriptorSource) (problems []lint.Problem, err error) {
				// We only care about Get methods for the purpose of this rule;
				// ignore everything else.
				if !isGetRequestMessage(m) {
					return
				}

				// Rule check: Establish that a name field is present.
				nameField := m.Fields().ByName("name")
				if nameField == nil {
					problems = append(problems, lint.Problem{
						Message:    fmt.Sprintf("method %q has no `name` field", m.Name()),
						Descriptor: m,
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
			MessageCallback: func(m p.MessageDescriptor, s lint.DescriptorSource) (problems []lint.Problem, err error) {
				// We only care about Get methods for the purpose of this rule;
				// ignore everything else.
				if !isGetRequestMessage(m) {
					return
				}

				// Rule check: Establish that there are no unexpected fields.
				allowedFields := map[string]struct{}{
					"name":      {}, // AIP-131
					"read_mask": {}, // AIP-157
					"view":      {}, // AIP-157
				}
				fields := m.Fields()
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
	return regexp.MustCompile("^Get(?:[A-Z]|$)").MatchString(methodName)
}

// Return true if this is an AIP-131 Get request message, false otherwise.
func isGetRequestMessage(m p.MessageDescriptor) bool {
	return regexp.MustCompile("^Get[A-Za-z0-9]*Request$").MatchString(string(m.Name()))
}
