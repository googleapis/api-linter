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
	"bitbucket.org/creachadair/stringset"
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

var httpBody = &lint.MethodRule{
	Name:   lint.NewRuleName("core", "0136", "http-body"),
	URI:    "https://aip.dev/136",
	OnlyIf: isCustomMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		for _, httpRule := range utils.GetHTTPRules(m) {
			noBody := stringset.New("GET", "DELETE")
			if !noBody.Contains(httpRule.Method) && httpRule.Body != "*" {
				// FIXME: We intentionally only return one Problem here.
				// When we can attach problems to the particular annotation, update this
				// to return multiples.
				return []lint.Problem{{
					Message:    "Custom POST methods should set `body: \"*\"`.",
					Descriptor: m,
				}}
			} else if noBody.Contains(httpRule.Method) && httpRule.Body != "" {
				// FIXME: We intentionally only return one Problem here.
				// When we can attach problems to the particular annotation, update this
				// to return multiples.
				return []lint.Problem{{
					Message:    "Custom GET (or DELETE) methods should not set a body clause.",
					Descriptor: m,
				}}
			}
		}
		return nil
	},
}
