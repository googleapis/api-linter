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

package aip0123

import (
	"strings"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

var resourceVariables = &lint.MessageRule{
	Name:   lint.NewRuleName(123, "resource-variables"),
	OnlyIf: hasResourceAnnotation,
	LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
		resource := utils.GetResource(m)
		for _, pattern := range resource.GetPattern() {
			for _, variable := range getVariables(pattern) {
				if strings.ToLower(variable) != variable {
					return []lint.Problem{{
						Message:    "Variable names in patterns should use snake case.",
						Descriptor: m,
						Location:   locations.MessageResource(m),
					}}
				}
				if strings.HasSuffix(variable, "_id") {
					return []lint.Problem{{
						Message:    "Variable names should omit the `_id` suffix.",
						Descriptor: m,
						Location:   locations.MessageResource(m),
					}}
				}
			}
		}
		return nil
	},
}
