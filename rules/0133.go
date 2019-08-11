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

	"github.com/stoewer/go-strcase"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/descriptor"
	p "google.golang.org/protobuf/reflect/protoreflect"
)

func init() {
	registerRules(
		checkCreateRequestMessageName(),
		checkCreateRequestMessageParentField(),
		checkCreateRequestMessageUnknownFields(),
		checkCreateResponseMessageName(),
	)
}

// Create messages should have a properly named Request message.
func checkCreateRequestMessageName() lint.Rule {
	return &descriptor.CallbackRule{
		RuleInfo: lint.RuleInfo{
			Name:         lint.NewRuleName("core", "0133", "request_message", "name"),
			Description:  "check that Create RPCs have appropriate request messages",
			RequestTypes: []lint.RequestType{lint.ProtoRequest},
			URI:          "https://aip.dev/133#guidance",
		},
		Callback: descriptor.Callbacks{
			MethodCallback: func(m p.MethodDescriptor, s lint.DescriptorSource) (problems []lint.Problem, err error) {
				// We only care about Create methods for the purpose of this rule;
				// ignore everything else.
				if !isCreateMethod(m) {
					return
				}

				// Rule check: Make sure there is an actual resource being created.
				// (The RPC should not just be "Create").
				methodName := string(m.Name())
				if methodName == "Create" {
					problems = append(problems, lint.Problem{
						Message:    "Create RPCs should indicate what they create (`CreateFoo`, not `Create`).",
						Descriptor: m,
					})
					return problems, nil
				}

				// Rule check: Establish that for methods such as `CreateFoo`, the request
				// message is named `CreateFooRequest`.
				requestMessageName := string(m.Input().Name())
				correctRequestMessageName := methodName + "Request"
				if requestMessageName != correctRequestMessageName {
					problems = append(problems, lint.Problem{
						Message: fmt.Sprintf(
							"Create RPCs should have a request message named after the RPC, such as %q.",
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

// The Create standard method should have a parent field.
func checkCreateRequestMessageParentField() lint.Rule {
	return &descriptor.CallbackRule{
		RuleInfo: lint.RuleInfo{
			Name:         lint.NewRuleName("core", "0133", "request_message", "parent_field"),
			Description:  "check that a parent field is present",
			RequestTypes: []lint.RequestType{lint.ProtoRequest},
			URI:          "https://aip.dev/133#request-message",
		},
		Callback: descriptor.Callbacks{
			MessageCallback: func(m p.MessageDescriptor, s lint.DescriptorSource) (problems []lint.Problem, err error) {
				// We only care about Create methods for the purpose of this rule;
				// ignore everything else.
				if !isCreateRequestMessage(m) {
					return
				}

				// Rule check: Establish that a parent field is present.
				parentField := m.Fields().ByName("parent")
				if parentField == nil {
					problems = append(problems, lint.Problem{
						Message:    fmt.Sprintf("method %q has no `parent` field", m.Name()),
						Descriptor: m,
					})
					return problems, nil
				}

				// Rule check: Establish that the parent field is a string.
				if parentField.Kind() != p.StringKind {
					problems = append(problems, lint.Problem{
						Message:    "The `parent` field on Create RPCs should be a string",
						Descriptor: parentField,
					})
				}

				return problems, nil
			},
		},
	}
}

// Create methods should not have unrecognized fields.
func checkCreateRequestMessageUnknownFields() lint.Rule {
	return &descriptor.CallbackRule{
		RuleInfo: lint.RuleInfo{
			Name:         lint.NewRuleName("core", "0133", "request_message", "unknown_fields"),
			Description:  "check that there are no unknown fields",
			RequestTypes: []lint.RequestType{lint.ProtoRequest},
			URI:          "https://aip.dev/133#request-message",
		},
		Callback: descriptor.Callbacks{
			MessageCallback: func(m p.MessageDescriptor, s lint.DescriptorSource) (problems []lint.Problem, err error) {
				// We only care about Get methods for the purpose of this rule;
				// ignore everything else.
				if !isCreateRequestMessage(m) {
					return
				}

				// Determine the name of the resource being created.
				resource := createRequestMessageRegexp.FindStringSubmatch(string(m.Name()))[1]
				resourceSnake := strcase.SnakeCase(resource)

				// Rule check: Establish that there are no unexpected fields.
				allowedFields := map[string]struct{}{
					"parent":                            {}, // AIP-133
					resourceSnake:                       {}, // AIP-133
					fmt.Sprintf("%s_id", resourceSnake): {}, // AIP-133
					"request_id":                        {}, // AIP-155
				}
				fields := m.Fields()
				for i := 0; i < fields.Len(); i++ {
					field := fields.Get(i)
					if _, ok := allowedFields[string(field.Name())]; !ok {
						problems = append(problems, lint.Problem{
							Message: fmt.Sprintf(
								"Create RPCs should only only contain fields explicitly described in AIPs, not %q.",
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

// Get messages should use the resource (or Operation) as the response message
func checkCreateResponseMessageName() lint.Rule {
	return &descriptor.CallbackRule{
		RuleInfo: lint.RuleInfo{
			Name:         lint.NewRuleName("core", "0133", "response_message", "name"),
			Description:  "check that Create RPCs have appropriate response messages",
			RequestTypes: []lint.RequestType{lint.ProtoRequest},
			URI:          "https://aip.dev/133#guidance",
		},
		Callback: descriptor.Callbacks{
			MethodCallback: func(m p.MethodDescriptor, s lint.DescriptorSource) (problems []lint.Problem, err error) {
				// We only care about Get methods for the purpose of this rule;
				// ignore everything else.
				if !isCreateMethod(m) {
					return
				}

				// Sanity check: If the response type is an LRO, accept this.
				// TODO: Check the response type annotation.
				if string(m.Output().FullName()) == "google.longrunning.Operation" {
					return problems, nil
				}

				// Rule check: Establish that for methods such as `CreateFoo`, the response
				// message is named `Foo`.
				responseMessageName := string(m.Output().Name())
				methodName := string(m.Name())
				if correctResponseMessageName := methodName[6:]; correctResponseMessageName != responseMessageName {
					problems = append(problems, lint.Problem{
						Message: fmt.Sprintf(
							"Create RPCs should have the corresponding resource as the response message, such as %q.",
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

var createMethodRegexp = regexp.MustCompile("^Create(?:[A-Z]|$)")
var createRequestMessageRegexp = regexp.MustCompile("^Create[A-Za-z0-9]+Request$")

// Return true if this is a AIP-133 Create method, false otherwise.
func isCreateMethod(m p.MethodDescriptor) bool {
	return createMethodRegexp.MatchString(string(m.Name()))
}

// Return true if this is an AIP-133 Create request message, false otherwise.
func isCreateRequestMessage(m p.MessageDescriptor) bool {
	return createRequestMessageRegexp.MatchString(string(m.Name()))
}
