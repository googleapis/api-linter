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

package aip0231

import (
	"fmt"

	"github.com/gertd/go-pluralize"
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

var pluralMethodResourceName = &lint.MethodRule{
	Name:   lint.NewRuleName("core", "0231", "method-name", "name"),
	OnlyIf: isBatchGetMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		// Note: m.GetName()[8:] is used to retrieve the resource name from the
		// method name. For example, "BatchGetFoos" -> "Foos"
		pluralMethodResourceName := m.GetName()[8:]

		pluralize := pluralize.NewClient()

		// Rule check: Establish that for methods such as `BatchGetFoos`
		if !pluralize.IsPlural(pluralMethodResourceName) {
			return []lint.Problem{{
				Message: fmt.Sprintf(
					`The resource part in method name %q shouldn't be %q, but should be its plural form %q`,
					m.GetName(), pluralMethodResourceName, pluralize.Plural(pluralMethodResourceName),
				),
				Descriptor: m,
			}}
		}

		return nil
	},
}

// Batch Get method should have a properly named Request message.
var inputName = &lint.MethodRule{
	Name:   lint.NewRuleName("core", "0231", "input-message", "name"),
	OnlyIf: isBatchGetMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		pluralInputResourceName := pluralize.NewClient().Plural(m.GetName()[8:])

		// Rule check: Establish that for methods such as `BatchGetFoos`, the request
		// message is named `BatchGetFoosRequest`.
		if got, want := m.GetInputType().GetName(), fmt.Sprintf("BatchGet%sRequest", pluralInputResourceName); got != want {
			return []lint.Problem{{
				Message: fmt.Sprintf(
					"Batch Get RPCs should have a properly named request message %q, but not %q",
					want, got,
				),
				Descriptor: m,
			}}
		}

		return nil
	},
}

// Batch Get method should have a properly named Response message.
var outputName = &lint.MethodRule{
	Name:   lint.NewRuleName("core", "0231", "output-message", "name"),
	OnlyIf: isBatchGetMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		pluralInputResourceName := pluralize.NewClient().Plural(m.GetName()[8:])

		// Rule check: Establish that for methods such as `BatchGetFoos`, the request
		// message is named `BatchGetFoosResponse`.
		if got, want := m.GetOutputType().GetName(), fmt.Sprintf("BatchGet%sResponse", pluralInputResourceName); got != want {
			return []lint.Problem{{
				Message: fmt.Sprintf(
					"Batch Get RPCs should have a properly named response message %q, but not %q",
					want, got,
				),
				Descriptor: m,
			}}
		}

		return nil
	},
}

// Batch Get methods should use the HTTP GET verb.
var httpVerb = &lint.MethodRule{
	Name:   lint.NewRuleName("core", "0231", "http-verb"),
	OnlyIf: isBatchGetMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		// Rule check: Establish that the RPC uses HTTP GET.
		for _, httpRule := range utils.GetHTTPRules(m) {
			if httpRule.Method != "GET" {
				return []lint.Problem{{
					Message:    "Batch Get methods must use the HTTP GET verb.",
					Descriptor: m,
				}}
			}
		}

		return nil
	},
}

// Batch Get methods should have a proper HTTP pattern.
var httpUrl = &lint.MethodRule{
	Name:   lint.NewRuleName("core", "0231", "http-name"),
	OnlyIf: isBatchGetMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		// Establish that the RPC has no HTTP body.
		for _, httpRule := range utils.GetHTTPRules(m) {
			if !batchGetURINameRegexp.MatchString(httpRule.URI) {
				return []lint.Problem{{
					Message:    `Batch Get methods URI should be end with ":batchGet".`,
					Descriptor: m,
				}}
			}
		}

		return nil
	},
}

// Get methods should not have an HTTP body.
var httpBody = &lint.MethodRule{
	Name:   lint.NewRuleName("core", "0231", "http-body"),
	OnlyIf: isBatchGetMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		// Establish that the RPC has no HTTP body.
		for _, httpRule := range utils.GetHTTPRules(m) {
			if httpRule.Body != "" {
				return []lint.Problem{{
					Message:    "Batch Get methods should not have an HTTP body.",
					Descriptor: m,
				}}
			}
		}

		return nil
	},
}
