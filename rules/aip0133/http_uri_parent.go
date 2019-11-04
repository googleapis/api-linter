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

package aip0133

import (
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

// Create methods should have a proper HTTP pattern.
var httpURIField = &lint.MethodRule{
	Name:   lint.NewRuleName("core", "0133", "http-uri-parent"),
	URI:    "https://aip.dev/133#guidance",
	OnlyIf: isCreateMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		// Establish that the RPC uri should include the `parent` field.
		for _, httpRule := range utils.GetHTTPRules(m) {
			if !createURINameRegexp.MatchString(httpRule.URI) {
				return []lint.Problem{{
					Message:    "Create methods should include the `parent` field in the URI.",
					Descriptor: m,
				}}
			}
		}

		return nil
	},
}
