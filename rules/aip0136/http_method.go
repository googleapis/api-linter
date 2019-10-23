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
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

var httpMethod = &lint.MethodRule{
	Name:   lint.NewRuleName("core", "0136", "http-method"),
	URI:    "https://aip.dev/136#guidance",
	OnlyIf: isCustomMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		for _, httpRule := range utils.GetHTTPRules(m) {
			if httpRule.Method != "POST" && httpRule.Method != "GET" {
				return []lint.Problem{{
					Message:    "Custom methods should use the HTTP POST or GET method.",
					Descriptor: m,
				}}
			}
		}
		return nil
	},
}
