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

package aip0131

import (
	"fmt"
	"strings"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

// Get messages should have a properly named Request message.
var requestMessageName = &lint.MethodRule{
	Name:   lint.NewRuleName("core", "0131", "request-message", "name"),
	URI:    "https://aip.dev/131#guidance",
	OnlyIf: isGetMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		// Rule check: Establish that for methods such as `GetFoo`, the request
		// message is named `GetFooRequest`.
		if got, want := m.GetInputType().GetName(), m.GetName()+"Request"; got != want {
			return []lint.Problem{{
				Message: fmt.Sprintf(
					"Get RPCs should have a request message named after the RPC, such as %q.",
					want,
				),
				Suggestion: want,
				Descriptor: m,
			}}
		}

		return nil
	},
}

// Get messages should use the resource as the response message
var responseMessageName = &lint.MethodRule{
	Name:   lint.NewRuleName("core", "0131", "response-message", "name"),
	URI:    "https://aip.dev/131#guidance",
	OnlyIf: isGetMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		// Rule check: Establish that for methods such as `GetFoo`, the response
		// message is named `Foo`.
		if got, want := m.GetOutputType().GetName(), m.GetName()[3:]; got != want {
			return []lint.Problem{{
				Message: fmt.Sprintf(
					"Get RPCs should have the corresponding resource as the response message, such as %q.",
					want,
				),
				Suggestion: want,
				Descriptor: m,
			}}
		}

		return nil
	},
}

// Get methods should use the HTTP GET verb.
var httpVerb = &lint.MethodRule{
	Name:   lint.NewRuleName("core", "0131", "http-verb"),
	URI:    "https://aip.dev/131#guidance",
	OnlyIf: isGetMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		// Rule check: Establish that the RPC uses HTTP GET.
		for _, httpRule := range utils.GetHTTPRules(m) {
			if httpRule.Method != "GET" {
				return []lint.Problem{{
					Message:    "Get methods must use the HTTP GET verb.",
					Descriptor: m,
				}}
			}
		}

		return nil
	},
}

// Get methods should have a proper HTTP pattern.
var httpNameField = &lint.MethodRule{
	Name:   lint.NewRuleName("core", "0131", "http-name"),
	URI:    "https://aip.dev/131#guidance",
	OnlyIf: isGetMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		// Establish that the RPC has no HTTP body.
		for _, httpRule := range utils.GetHTTPRules(m) {
			if !getURINameRegexp.MatchString(httpRule.URI) {
				return []lint.Problem{{
					Message:    "Get methods should include the `name` field in the URI.",
					Descriptor: m,
				}}
			}
		}

		return nil
	},
}

// Get methods should not have an HTTP body.
var httpBody = &lint.MethodRule{
	Name:   lint.NewRuleName("core", "0131", "http-body"),
	URI:    "https://aip.dev/131#guidance",
	OnlyIf: isGetMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		// Establish that the RPC has no HTTP body.
		for _, httpRule := range utils.GetHTTPRules(m) {
			if httpRule.Body != "" {
				return []lint.Problem{{
					Message:    "Get methods should not have an HTTP body.",
					Descriptor: m,
				}}
			}
		}

		return nil
	},
}

// Get methods should not generally use synonyms for "get".
var synonyms = &lint.MethodRule{
	Name: lint.NewRuleName("core", "0131", "synonyms"),
	URI:  "https://aip.dev/131#guidance",
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		name := m.GetName()
		for _, syn := range []string{"Acquire", "Fetch", "Lookup", "Read", "Retrieve"} {
			if strings.HasPrefix(name, syn) {
				return []lint.Problem{{
					Message: fmt.Sprintf(
						`%q can be a synonym for "Get". Should this be a Get method?`,
						syn,
					),
					Suggestion: strings.Replace(name, syn, "Get", 1),
					Descriptor: m,
					Location:   lint.DescriptorNameLocation(m),
				}}
			}
		}
		return nil
	},
}
