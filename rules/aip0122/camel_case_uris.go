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
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

// HTTP URL pattern shouldn't include underscore("_")
var httpURICase = &lint.MethodRule{
	Name: lint.NewRuleName("core", "0122", "camel-case-uris"),
	URI:  "https://aip.dev/122#guidance",
	LintMethod: func(m *desc.MethodDescriptor) (problems []lint.Problem) {
		// Establish that the URI does not include a `_` character.
		for _, httpRule := range utils.GetHTTPRules(m) {
			parsedURI := parseURI(httpRule.URI)
			if strings.Contains(parsedURI.pattern, "_") {
				problems = append(problems, lint.Problem{
					Message:    "HTTP URI patterns should use camel case, not snake case.",
					Descriptor: m,
				})
			}
			for _, v := range parsedURI.vars {
				if strings.ToLower(v) != v {
					problems = append(problems, lint.Problem{
						Message:    "Variable names in URI patterns should use snake case, not camel case.",
						Descriptor: m,
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

// extractURI separates variable names from the URI pattern.
//
// It does not maintain the relationship between variable names and what
// part of the pattern they apply to (we do not need that to lint here).
func parseURI(uri string) (parsed struct {
	vars    []string
	pattern string
}) {
	parsed.pattern = uri
	for strings.Contains(parsed.pattern, "{") && strings.Contains(parsed.pattern, "}") {
		// Find the first {variable=pattern} segment and pull the variable out of it.
		start, end := strings.Index(uri, "{"), strings.Index(uri, "}")
		repl := ""
		for i, segment := range strings.SplitN(parsed.pattern[start+1:end], "=", 2) {
			if i == 0 {
				parsed.vars = append(parsed.vars, segment)
			} else {
				repl = segment
			}
		}
		parsed.pattern = strings.Replace(parsed.pattern, parsed.pattern[start:end+1], repl, 1)
	}
	return
}
