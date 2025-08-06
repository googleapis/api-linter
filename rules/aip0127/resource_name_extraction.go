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

package aip0127

import (
	"github.com/googleapis/api-linter/v2/lint"
	"github.com/googleapis/api-linter/v2/locations"
	"github.com/googleapis/api-linter/v2/rules/internal/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var resourceNameExtraction = &lint.MethodRule{
	Name: lint.NewRuleName(127, "resource-name-extraction"),
	LintMethod: func(m protoreflect.MethodDescriptor) []lint.Problem {
		for _, rule := range utils.GetHTTPRules(m) {
			for k, v := range rule.GetVariables() {
				if v == "*" && k != "$api_version" {
					return []lint.Problem{{
						Message:    "Extract a full resource name into a variable, not just IDs.",
						Descriptor: m,
						Location:   locations.MethodHTTPRule(m),
					}}
				}
			}
		}
		return nil
	},
}
