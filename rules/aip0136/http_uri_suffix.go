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

package aip0136

import (
	"fmt"
	"strings"

	"github.com/googleapis/api-linter/locations"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
	"github.com/stoewer/go-strcase"
)

var uriSuffix = &lint.MethodRule{
	Name: lint.NewRuleName(136, "http-uri-suffix"),
	OnlyIf: func(m *desc.MethodDescriptor) bool {
		return isCustomMethod(m) && httpNameVariable.LintMethod(m) == nil && httpParentVariable.LintMethod(m) == nil
	},
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		for _, httpRule := range utils.GetHTTPRules(m) {
			// URIs should end in a `:` character followed by the name of the method.
			// However, if the noun is the URI's resource, then we only use `:verb`,
			// not `:verbNoun`.
			//
			// This is somewhat tricky to test for perfectly, and may need to evolve
			// over time, but the following rules should be mostly correct:
			//   1. If the URI contains `{name=` or `{parent=`, expect `:verb`.
			//   2. Otherwise, expect `:verbNoun`.
			//
			// We blindly assume that the verb is always one word (the "noun" may be
			// any number of words; they often include adjectives).
			//
			// N.B. The LowerCamel(Snake(name)) is because strcase does not translate
			//      from upper camel to lower camel correctly.
			want := ":" + strcase.LowerCamelCase(strcase.SnakeCase(m.GetName()))
			if strings.Contains(httpRule.URI, "{name=") || strings.Contains(httpRule.URI, "{parent=") {
				rpcSlice := strings.Split(strcase.SnakeCase(m.GetName()), "_")
				if rpcSlice[0] == "batch" {
					want = ":" + strcase.LowerCamelCase(rpcSlice[0]+"_"+rpcSlice[1])
				} else {
					want = ":" + rpcSlice[0]
				}
			}

			// AIP-162 introduces some special cases around revisions, where
			// `ListFooRevisions` gets a suffix of `:listRevisions` (and the same for
			// `Delete` and `Tag`).
			n := m.GetName()
			if strings.HasPrefix(n, "List") && strings.HasSuffix(m.GetName(), "Revisions") {
				want = ":listRevisions"
			}
			if strings.HasSuffix(m.GetName(), "Revision") {
				if strings.HasPrefix(m.GetName(), "Tag") {
					want = ":tagRevision"
				}
				if strings.HasPrefix(m.GetName(), "Delete") {
					want = ":deleteRevision"
				}
			}

			// Do we have the suffix we expect?
			if !strings.HasSuffix(httpRule.URI, want) {
				// FIXME: We intentionally only return one Problem here.
				// When we can attach issues to the particular annotation, update this
				// to return multiples.
				return []lint.Problem{{
					Message: fmt.Sprintf(
						"Custom method should have a URI suffix matching the method name, such as %q.",
						want,
					),
					Descriptor: m,
					Location:   locations.MethodHTTPRule(m),
				}}
			}
		}
		return nil
	},
}
