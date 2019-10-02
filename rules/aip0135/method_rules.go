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

package aip0135

import (
	"fmt"

	"github.com/googleapis/api-linter/internal/utils"
	"github.com/googleapis/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
)

// Delete messages should have a properly named Request message.
var requestMessageName = &lint.MethodRule{
	Name:   lint.NewRuleName("core", "0135", "request-message", "name"),
	URI:    "https://aip.dev/135#guidance",
	OnlyIf: isDeleteMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		// Rule check: Establish that for methods such as `DeleteFoo`, the request
		// message is named `DeleteFooRequest`.
		if got, want := m.GetInputType().GetName(), m.GetName()+"Request"; got != want {
			return []lint.Problem{{
				Message: fmt.Sprintf(
					"Delete RPCs should have a request message named after the RPC, such as %q.",
					want,
				),
				Suggestion: want,
				Descriptor: m,
			}}
		}

		return nil
	},
}

// Delete messages should use google.protobuf.Empty, google.longrunning.Operation or the resource itself as the response message
var responseMessageName = &lint.MethodRule{
	Name:   lint.NewRuleName("core", "0135", "response-message", "name"),
	URI:    "https://aip.dev/135#guidance",
	OnlyIf: isDeleteMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		// Rule check: Establish that for methods such as `DeleteFoo`, the response
		// message is `google.protobuf.Empty`, `google.longrunning.Operation`, or `Foo`.
		got := m.GetOutputType().GetName()
		want := []string{"Empty", "Operation", m.GetName()[6:]}
		for _, v := range want {
			if got == v {
				return nil
			}
		}
		return []lint.Problem{{
			Message: fmt.Sprintf(
				"Delete RPCs should have response message type in %q, not %q.",
				want,
				got,
			),
			Suggestion: want[0],
			Descriptor: m,
		}}
	},
}

// Delete methods should use the HTTP DELETE verb.
var httpVerb = &lint.MethodRule{
	Name:   lint.NewRuleName("core", "0135", "http-verb"),
	URI:    "https://aip.dev/135#guidance",
	OnlyIf: isDeleteMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		// Rule check: Establish that the RPC uses HTTP DELETE.
		for _, httpRule := range utils.GetHTTPRules(m) {
			if httpRule.GetDelete() == "" {
				return []lint.Problem{{
					Message:    "Delete methods must use the HTTP DELETE verb.",
					Descriptor: m,
				}}
			}
		}

		return nil
	},
}

// Delete methods should have a proper HTTP pattern.
var httpNameField = &lint.MethodRule{
	Name:   lint.NewRuleName("core", "0135", "http-name"),
	URI:    "https://aip.dev/135#guidance",
	OnlyIf: isDeleteMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		// Establish that the RPC has no HTTP body.
		for _, httpRule := range utils.GetHTTPRules(m) {
			if uri := httpRule.GetDelete(); uri != "" {
				if !deleteURINameRegexp.MatchString(uri) {
					return []lint.Problem{{
						Message:    "Delete methods should include the `name` field in the URI.",
						Descriptor: m,
					}}
				}
			}
		}

		return nil
	},
}

// Delete methods should not have an HTTP body.
var httpBody = &lint.MethodRule{
	Name:   lint.NewRuleName("core", "0135", "http-body"),
	URI:    "https://aip.dev/131#guidance",
	OnlyIf: isDeleteMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		// Establish that the RPC has no HTTP body.
		for _, httpRule := range utils.GetHTTPRules(m) {
			if httpRule.GetBody() != "" {
				return []lint.Problem{{
					Message:    "Delete methods should not have an HTTP body.",
					Descriptor: m,
				}}
			}
		}

		return nil
	},
}
