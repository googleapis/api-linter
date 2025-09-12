// Copyright 2025 Google LLC
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

package aip0160

import (
	"fmt"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

const (
	typeMessageFmt = `Guidance: use the field, "string filter", not "%s filter"`
	typeSuggestion = "string"
)

var filterFieldType = &lint.MethodRule{
	Name: lint.NewRuleName(160, "filter-field-type"),
	OnlyIf: func(m *desc.MethodDescriptor) bool {
		return utils.IsListMethod(m) || utils.IsCustomMethod(m)
	},
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		var problems []lint.Problem
		for _, f := range m.GetInputType().GetFields() {
			if f.GetName() == "filter" && utils.GetTypeName(f) != "string" {
				problems = append(problems, lint.Problem{
					Message:    fmt.Sprintf(typeMessageFmt, utils.GetTypeName(f)),
					Descriptor: f,
					Location:   locations.FieldType(f),
					Suggestion: typeSuggestion,
				})
			}
		}
		return problems
	},
}
