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

package aip0134

import (
	"fmt"
	"strings"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
	"github.com/stoewer/go-strcase"
)

// Update methods should have a properly named Request message.
var requestMessageName = &lint.MethodRule{
	Name:   lint.NewRuleName("core", "0134", "request-message", "name"),
	URI:    "https://aip.dev/134#guidance",
	OnlyIf: isUpdateMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		// Rule check: Establish that for methods such as `UpdateFoo`, the request
		// message is named `UpdateFooRequest`.
		if got, want := m.GetInputType().GetName(), m.GetName()+"Request"; got != want {
			return []lint.Problem{{
				Message: fmt.Sprintf(
					"Update RPCs should have a request message named after the RPC, such as %q.",
					want,
				),
				Suggestion: want,
				Descriptor: m,
			}}
		}

		return nil
	},
}

// Update methods should use the resource as the response message
var responseMessageName = &lint.MethodRule{
	Name:   lint.NewRuleName("core", "0134", "response-message", "name"),
	URI:    "https://aip.dev/134#guidance",
	OnlyIf: isUpdateMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		// Rule check: Establish that for methods such as `UpdateFoo`, the response
		// message is `Foo` or `google.longrunning.Operation`.
		got := m.GetOutputType().GetFullyQualifiedName()
		want := []string{
			m.GetName()[6:],
			"google.longrunning.Operation",
		}
		for _, v := range want {
			if got == v {
				return nil
			}
		}
		return []lint.Problem{{
			Message: fmt.Sprintf(
				"Update RPCs should have response message type in %q, not %q.",
				want,
				got,
			),
			Suggestion: want[0],
			Descriptor: m,
		}}
	},
}

// Update methods should use the HTTP PATCH verb.
var httpVerb = &lint.MethodRule{
	Name:   lint.NewRuleName("core", "0134", "http-verb"),
	URI:    "https://aip.dev/134#patch-and-put",
	OnlyIf: isUpdateMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		// Rule check: Establish that the RPC uses HTTP PATCH.
		for _, httpRule := range utils.GetHTTPRules(m) {
			if httpRule.GetPatch() == "" {
				return []lint.Problem{{
					Message:    "Update methods must use the HTTP PATCH verb.",
					Descriptor: m,
				}}
			}
		}

		return nil
	},
}

// Update methods should have a proper HTTP pattern.
var httpNameField = &lint.MethodRule{
	Name:   lint.NewRuleName("core", "0134", "http-name"),
	URI:    "https://aip.dev/134#guidance",
	OnlyIf: isUpdateMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		fieldName := strcase.SnakeCase(m.GetName()[6:])
		// Establish that the RPC has expected HTTP pattern.
		for _, httpRule := range utils.GetHTTPRules(m) {
			if uri := httpRule.GetPatch(); uri != "" {
				matches := updateURINameRegexp.FindStringSubmatch(uri)
				if matches == nil || matches[1] != fieldName {
					return []lint.Problem{{
						Message:    fmt.Sprintf("Update methods should include the `%s.name` field in the URI.", fieldName),
						Descriptor: m,
					}}
				}
			}
		}

		return nil
	},
}

// Update methods should have an HTTP body.
var httpBody = &lint.MethodRule{
	Name:   lint.NewRuleName("core", "0134", "http-body"),
	URI:    "https://aip.dev/134#guidance",
	OnlyIf: isUpdateMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		fieldName := strcase.SnakeCase(m.GetName()[6:])
		// Establish that the RPC has HTTP body equal to fieldName.
		for _, httpRule := range utils.GetHTTPRules(m) {
			if httpRule.GetBody() != fieldName {
				return []lint.Problem{{
					Message:    fmt.Sprintf("Update methods should have an HTTP body equal to `%q`.", fieldName),
					Descriptor: m,
				}}
			}
		}

		return nil
	},
}

// Update methods should use the word "update", not synonyms.
var synonyms = &lint.MethodRule{
	Name: lint.NewRuleName("core", "0134", "synonyms"),
	URI:  "https://aip.dev/134#guidance",
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		name := m.GetName()
		for _, syn := range []string{"Patch", "Put", "Set"} {
			if strings.HasPrefix(name, syn) {
				return []lint.Problem{{
					Message: fmt.Sprintf(
						`%q can be a synonym for "Update". Should this be a Update method?`,
						syn,
					),
					Descriptor: m,
					Location:   lint.DescriptorNameLocation(m),
					Suggestion: strings.Replace(name, syn, "Update", 1),
				}}
			}
		}
		return nil
	},
}
