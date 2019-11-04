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
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

// Delete methods should use the HTTP DELETE method.
var httpMethod = &lint.MethodRule{
	Name:   lint.NewRuleName("core", "0135", "http-method"),
	OnlyIf: isDeleteMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		// Rule check: Establish that the RPC uses HTTP DELETE.
		for _, httpRule := range utils.GetHTTPRules(m) {
			if httpRule.Method != "DELETE" {
				return []lint.Problem{{
					Message:    "Delete methods must use the HTTP DELETE verb.",
					Descriptor: m,
				}}
			}
		}

		return nil
	},
}
