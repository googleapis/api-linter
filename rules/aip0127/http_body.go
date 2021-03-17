// Copyright 2021 Google LLC
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

package aip0127

import (
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

// GET and DELETE methods must not have an HTTP body.
var httpBody = &lint.MethodRule{
	Name: lint.NewRuleName(127, "http-body"),
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		var ps []lint.Problem
		for _, r := range utils.GetHTTPRules(m) {
			if r.Method != "GET" && r.Method != "DELETE" {
				continue
			}
			if r.Body != "" {
				ps = append(ps, lint.Problem{
					Message:    "`GET` and `DELETE` HTTP methods must not have an HTTP `body`.",
					Descriptor: m,
					Location:   locations.MethodHTTPRule(m),
				})
			}
		}
		return ps
	},
}
