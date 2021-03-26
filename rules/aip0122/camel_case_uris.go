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

package aip0122

import (
	"strings"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

// HTTP URL pattern shouldn't include underscore("_")
var httpURICase = &lint.MethodRule{
	Name: lint.NewRuleName(122, "camel-case-uris"),
	LintMethod: func(m *desc.MethodDescriptor) (problems []lint.Problem) {
		// Establish that the URI does not include a `_` character.
		for _, httpRule := range utils.GetHTTPRules(m) {
			if strings.Contains(httpRule.GetPlainURI(), "_") {
				problems = append(problems, lint.Problem{
					Message:    "HTTP URI patterns should use camel case, not snake case.",
					Descriptor: m,
					Location:   locations.MethodHTTPRule(m),
				})
			}
			for v := range httpRule.GetVariables() {
				if strings.ToLower(v) != v {
					problems = append(problems, lint.Problem{
						Message:    "Variable names in URI patterns should use snake case, not camel case.",
						Descriptor: m,
						Location:   locations.MethodHTTPRule(m),
					})
				}
			}

			// FIXME: We intentionally only return at most one of each type of `Problem` here.
			// When we can attach problems to the particular annotation, remove this.
			if len(problems) > 0 {
				return
			}
		}
		return
	},
}
