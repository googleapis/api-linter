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
var camelCase = &lint.MethodRule{
	Name: lint.NewRuleName("core", "0122", "http-url-pattern"),
	URI:  "https://aip.dev/122#guidance",
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		// Establish that the URI does not include a `_` character.
		for _, httpRule := range utils.GetHTTPRules(m) {
			if strings.Contains(httpRule.URI, "_") {
				return []lint.Problem{{
					Message:    "HTTP URL should use camel case, but not snake case.",
					Descriptor: m,
				}}
			}
		}
		return nil
	},
}
