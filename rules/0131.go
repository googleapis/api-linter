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
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"
)

func init() {
	registerRules(
		checkGetRequestMessageName,
		checkGetRequestMessageNameField,
		checkGetRequestMessageUnknownFields,
		checkGetResponseMessageName,
	)
}

// Get messages should have a properly named Request message.
var checkGetRequestMessageName = lint.Rule{
	Name:         lint.NewRuleName("core", "0131", "request-message", "name"),
	Description:  "Get RPCs must have a consistent request message name.",
	URI:          "https://aip.dev/131#guidance",
	LintMethod: func(m *desc.MethodDescriptor) (problems []lint.Problem) {
				// We only care about Get methods for the purpose of this rule;
				// ignore everything else.
				if !isGetMethod(m) {
					return problems
				}

				// Rule check: Establish that for methods such as `GetFoo`, the request
				// message is named `GetFooRequest`.
				if got, want := m.GetInputType().GetName(), m.GetName() + "Request"; got != want {
					problems = append(problems, lint.Problem{
						Message: fmt.Sprintf(
							"Get RPCs should have a request message named after the RPC, such as %q.",
							want,
						),
						Suggestion: want,
						Descriptor: m,
					})
				}

				return problems
			},
}

// The Get standard method should only have expected fields.
var checkGetRequestMessageNameField = lint.Rule{
	Name:         lint.NewRuleName("core", "0131", "request-message", "name-field"),
	Description:  "check that a name field is present",
	URI:          "https://aip.dev/131#request-message",
	LintMessage: func(m *desc.MessageDescriptor) (problems []lint.Problem) {
				// We only care about Get methods for the purpose of this rule;
				// ignore everything else.
				if !isGetRequestMessage(m) {
					return problems
				}

				// Rule check: Establish that a name field is present.
				name := m.FindFieldByName("name")
				if name == nil {
					problems = append(problems, lint.Problem{
						Message:    fmt.Sprintf("method %q has no `name` field", m.GetName()),
						Descriptor: m,
					})
					return problems
				}

				// Rule check: Ensure that the name field is the correct type.
				if name.GetType() != builder.FieldTypeString().GetType() {
					problems = append(problems, lint.Problem{
						Message:    "`name` field on Get RPCs should be a string",
						Descriptor: name,
					})
				}

				return problems
			},
}

// Get methods should not have unrecognized fields.
var checkGetRequestMessageUnknownFields = lint.Rule {
	Name:         lint.NewRuleName("core", "0131", "request-message", "unknown-fields"),
	Description:  "Get RPCs must not contain unexpected fields.",
	URI:          "https://aip.dev/131#request-message",
	LintMessage: func (m *desc.MessageDescriptor) (problems []lint.Problem) {
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
				for _, field := range m.GetFields() {
					if _, ok := allowedFields[string(field.GetName())]; !ok {
						problems = append(problems, lint.Problem{
							Message: fmt.Sprintf(
								"Get RPCs must only contain fields explicitly described in AIPs, not %q.",
								string(field.GetName()),
							),
							Descriptor: field,
						})
					}
				}

				return problems
			},
}

// Get messages should use the resource as the response message
var checkGetResponseMessageName = lint.Rule {
	Name:         lint.NewRuleName("core", "0131", "response-message", "name"),
	Description:  "check that Get RPCs have appropriate response messages",
	URI:          "https://aip.dev/131#guidance",
	LintMethod: func (m *desc.MethodDescriptor) (problems []lint.Problem) {
				// We only care about Get methods for the purpose of this rule;
				// ignore everything else.
				if !isGetMethod(m) {
					return
				}

				// Rule check: Establish that for methods such as `GetFoo`, the response
				// message is named `Foo`.
				if got, want := m.GetOutputType().GetName(), m.GetName()[3:]; got != want {
					problems = append(problems, lint.Problem{
						Message: fmt.Sprintf(
							"Get RPCs should have the corresponding resource as the response message, such as %q.",
							want,
						),
						Suggestion: want,
						Descriptor: m,
					})
				}

				return problems
			},
}

// Return true if this is a AIP-131 Get method, false otherwise.
func isGetMethod(m *desc.MethodDescriptor) bool {
	methodName := m.GetName()
	if methodName == "GetIamPolicy" {
		return false
	}
	return regexp.MustCompile("^Get(?:[A-Z]|$)").MatchString(methodName)
}

// Return true if this is an AIP-131 Get request message, false otherwise.
func isGetRequestMessage(m *desc.MessageDescriptor) bool {
	return regexp.MustCompile("^Get[A-Za-z0-9]*Request$").MatchString(m.GetName())
}
