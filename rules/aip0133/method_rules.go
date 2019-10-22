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
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

// Create methods should use the HTTP POST verb.
var httpVerb = &lint.MethodRule{
	Name:   lint.NewRuleName("core", "0133", "http-verb"),
	URI:    "https://aip.dev/133#guidance",
	OnlyIf: isCreateMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		// Rule check: Establish that the RPC uses HTTP POST.
		for _, httpRule := range utils.GetHTTPRules(m) {
			if httpRule.Method != "POST" {
				return []lint.Problem{{
					Message:    "Create methods must use the HTTP POST verb.",
					Descriptor: m,
				}}
			}
		}

		return nil
	},
}

// Create methods should have a proper HTTP pattern.
var httpURIField = &lint.MethodRule{
	Name:   lint.NewRuleName("core", "0133", "http-uri"),
	URI:    "https://aip.dev/133#guidance",
	OnlyIf: isCreateMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		// Establish that the RPC uri should include the `parent` field.
		for _, httpRule := range utils.GetHTTPRules(m) {
			if !createURINameRegexp.MatchString(httpRule.URI) {
				return []lint.Problem{{
					Message:    "Create methods should include the `parent` field in the URI.",
					Descriptor: m,
				}}
			}
		}

		return nil
	},
}

// Create methods should have an HTTP body, and the body value should be resource.
var httpBody = &lint.MethodRule{
	Name:   lint.NewRuleName("core", "0133", "http-body"),
	URI:    "https://aip.dev/133#guidance",
	OnlyIf: isCreateMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		// We only care about Create methods for the purpose of this rule;
		// ignore everything else.
		if !isCreateMethod(m) {
			return nil
		}

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
				}}
				// When resource field is not set in the request message, the problem
				// will not be triggered by the rule"core::0133::http-body". It will be
				// triggered by another rule
				// "core::0133::request-message::resource-field"
			} else if resourceFieldName != "" && httpRule.Body != resourceFieldName {
				return []lint.Problem{{
					Message:    fmt.Sprintf("The content of body %q must map to the resource field %q in the request message", body, resourceFieldName),
					Descriptor: m,
				}}
			}
		}
		return nil
	},
}

// Create method should have a properly named input message.
var inputName = &lint.MethodRule{
	Name:   lint.NewRuleName("core", "0133", "request-message", "name"),
	URI:    "https://aip.dev/133#guidance",
	OnlyIf: isCreateMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		resourceMsgName := getResourceMsgName(m)

		// Rule check: Establish that for methods such as `CreateFoo`, the request
		// message is named `CreateFooRequest`.
		if got, want := m.GetInputType().GetName(), fmt.Sprintf("Create%sRequest", resourceMsgName); got != want {
			return []lint.Problem{{
				Message: fmt.Sprintf(
					"Post RPCs should have a properly named request message %q, but not %q",
					want, got,
				),
				// TODO: suggestion will be set after the location is set properly
				// Suggestion: want,
				Descriptor: m,
			}}
		}

		return nil
	},
}

// Create method should use the resource as the output message
var outputName = &lint.MethodRule{
	Name:   lint.NewRuleName("core", "0133", "response-message", "name"),
	URI:    "https://aip.dev/133#guidance",
	OnlyIf: isCreateMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		output := m.GetOutputType()

		// AIP-0151
		if output.GetFullyQualifiedName() == "google.longrunning.Operation" {
			return nil
		}

		resourceMsgName := getResourceMsgName(m)

		// Rule check: Establish that for methods such as `CreateFoo`, the response
		// message should be named `Foo`
		if output.GetName() != resourceMsgName {
			return []lint.Problem{{
				Message: fmt.Sprintf(
					"Create RPCs should have the corresponding resource as the response message, such as %q.",
					resourceMsgName,
				),
				// TODO: suggestion will be set after the location is set properly
				// Suggestion: want,
				Descriptor: m,
			}}
		}

		return nil
	},
}

// Create methods should use "create", not synonyms.
var synonyms = &lint.MethodRule{
	Name: lint.NewRuleName("core", "0133", "synonyms"),
	URI:  "https://aip.dev/133",
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		name := m.GetName()
		for _, syn := range []string{"Insert", "Make", "Post"} {
			if strings.HasPrefix(name, syn) {
				return []lint.Problem{{
					Message: fmt.Sprintf(
						`%q can be a synonym for "Create". Should this be a Create method?`,
						syn,
					),
					Descriptor: m,
					Location:   lint.DescriptorNameLocation(m),
					Suggestion: strings.Replace(name, syn, "Create", 1),
				}}
			}
		}
		return nil
	},
}
