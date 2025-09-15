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

	"github.com/googleapis/api-linter/v2/lint"
	"github.com/googleapis/api-linter/v2/locations"
	"github.com/googleapis/api-linter/v2/rules/internal/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	nameMessageFmt = `Guidance: use the field, "string filter", not "string %s"`
	nameSuggestion = "filter"
)

var filterFieldName = &lint.MethodRule{
	Name: lint.NewRuleName(160, "filter-field-name"),
	OnlyIf: func(m protoreflect.MethodDescriptor) bool {
		return utils.IsListMethod(m) || utils.IsCustomMethod(m)
	},
	LintMethod: func(m protoreflect.MethodDescriptor) []lint.Problem {
		var problems []lint.Problem
		for ndx := 0; ndx < m.Input().Fields().Len(); ndx++ {
			f := m.Input().Fields().Get(ndx)
			if f.Name() == "filters" {
				problems = append(problems, lint.Problem{
					Message:    fmt.Sprintf(nameMessageFmt, f.Name()),
					Descriptor: f,
					Location:   locations.DescriptorName(f),
					Suggestion: nameSuggestion,
				})
			}
		}
		return problems
	},
}
