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

package aip0233

import (
	"fmt"

	"github.com/gertd/go-pluralize"
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

var pluralMethodResourceName = &lint.MethodRule{
	Name:   lint.NewRuleName("core", "0233", "method-name", "name"),
	URI:    "https://aip.dev/233#guidance",
	OnlyIf: isBatchCreateMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		// Note: m.GetName()[11:] is used to retrieve the resource name from the
		// method name. For example, "BatchCreateFoos" -> "Foos"
		pluralMethodResourceName := m.GetName()[11:]

		pluralize := pluralize.NewClient()

		// Rule check: Establish that for methods such as `BatchCreateFoos`
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

// Batch Create method should have a properly named Request message.
var inputName = &lint.MethodRule{
	Name:   lint.NewRuleName("core", "0233", "input-message", "name"),
	URI:    "https://aip.dev/233#guidance",
	OnlyIf: isBatchCreateMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		pluralInputResourceName := pluralize.NewClient().Plural(m.GetName()[11:])

		// Rule check: Establish that for methods such as `BatchCreateFoos`, the request
		// message is named `BatchCreateFoosRequest`.
		if got, want := m.GetInputType().GetName(), fmt.Sprintf("BatchCreate%sRequest", pluralInputResourceName); got != want {
			return []lint.Problem{{
				Message: fmt.Sprintf(
					"Batch Create RPCs should have a properly named request message %q, but not %q",
					want, got,
				),
				Descriptor: m,
			}}
		}

		return nil
	},
}

// Batch Create method should have a properly named Response message.
var outputName = &lint.MethodRule{
	Name:   lint.NewRuleName("core", "0233", "output-message", "name"),
	URI:    "https://aip.dev/233#guidance",
	OnlyIf: isBatchCreateMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		want := fmt.Sprintf("BatchCreate%sResponse", pluralize.NewClient().Plural(m.GetName()[11:]))

		// If this is an LRO, then use the annotated response type instead of
		// the actual RPC return type.
		got := m.GetOutputType().GetName()
		if m.GetOutputType().GetFullyQualifiedName() == "google.longrunning.Operation" {
			got = utils.GetOperationInfo(m).GetResponseType()
		}

		// Rule check: Establish that for methods such as `BatchCreateFoos`, the response
		// message should be named `BatchCreateFoosResponse`
		//
		// Note: If `got` is empty string, this is an unannotated LRO.
		// The AIP-151 rule will whine about that, and this rule should not as it
		// would be confusing.
		if got != want && got != "" {
			return []lint.Problem{{
				Message: fmt.Sprintf(
					"Batch Create RPCs should have a properly named response message %q, but not %q",
					want, got,
				),
				Descriptor: m,
			}}
		}

		return nil
	},
}

// Batch Create methods should use the HTTP POST verb.
var httpVerb = &lint.MethodRule{
	Name:   lint.NewRuleName("core", "0233", "http-verb"),
	URI:    "https://aip.dev/233#guidance",
	OnlyIf: isBatchCreateMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		// Rule check: Establish that the RPC uses HTTP POST.
		for _, httpRule := range utils.GetHTTPRules(m) {
			if httpRule.Method != "POST" {
				return []lint.Problem{{
					Message:    "Batch Create methods must use the HTTP POST verb.",
					Descriptor: m,
				}}
			}
		}

		return nil
	},
}

// Batch Create methods should have a proper HTTP pattern.
var httpUrl = &lint.MethodRule{
	Name:   lint.NewRuleName("core", "0233", "http-name"),
	URI:    "https://aip.dev/233#guidance",
	OnlyIf: isBatchCreateMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		// Establish that the RPC has no HTTP body.
		for _, httpRule := range utils.GetHTTPRules(m) {
			if !batchCreateURINameRegexp.MatchString(httpRule.URI) {
				return []lint.Problem{{
					Message:    `Batch Create methods URI should be end with ":batchCreate".`,
					Descriptor: m,
				}}
			}
		}

		return nil
	},
}

// Batch create methods should use "*" as the HTTP body.
var httpBody = &lint.MethodRule{
	Name:   lint.NewRuleName("core", "0233", "http-body"),
	URI:    "https://aip.dev/233#guidance",
	OnlyIf: isBatchCreateMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		// Establish that the RPC has no HTTP body.
		for _, httpRule := range utils.GetHTTPRules(m) {
			if httpRule.Body != "*" {
				return []lint.Problem{{
					Message:    `Batch Create methods should use "*" as the HTTP body.`,
					Descriptor: m,
				}}
			}
		}

		return nil
	},
}
