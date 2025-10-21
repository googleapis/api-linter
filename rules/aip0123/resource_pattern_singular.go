// Copyright 2024 Google LLC
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
	"fmt"
	"strings"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
	"github.com/stoewer/go-strcase"
)

var resourcePatternSingular = &lint.MessageRule{
	Name: lint.NewRuleName(123, "resource-pattern-singular"),
	OnlyIf: func(m *desc.MessageDescriptor) bool {
		return utils.IsResource(m) && len(utils.GetResource(m).GetPattern()) > 0
	},
	LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
		var problems []lint.Problem
		res := utils.GetResource(m)
		nested := isNestedName(res)
		var nn string
		if nested {
			nn = nestedSingular(res)
		}

		patterns := res.GetPattern()
		singular := utils.GetResourceSingular(res)

		if !utils.IsSingletonResource(m) {
			singular = fmt.Sprintf("{%s}", strcase.SnakeCase(singular))
			nn = fmt.Sprintf("{%s}", nn)
		} else {
			// singular is already in lower camel case
			// but nested name is returned in snake_case form
			// and final segment of singleton is lowerCamelCase
			nn = strcase.LowerCamelCase(nn)
		}

		// If the first pattern is reduced or non-compliant, but is nested name eligible, we want to recommend the nested name.
		nestedFirstPattern := nested && (strings.HasSuffix(patterns[0], nn) || !strings.HasSuffix(patterns[0], singular))

		for _, pattern := range patterns {
			if !strings.HasSuffix(pattern, singular) {
				// allow the reduced, nested name instead if present
				if nested && strings.HasSuffix(pattern, nn) {
					continue
				}
				want := singular
				if nestedFirstPattern {
					want = nn
				}

				problems = append(problems, lint.Problem{
					Message:    fmt.Sprintf("resource pattern %q final segment must include the resource singular %q", pattern, want),
					Descriptor: m,
					Location:   locations.MessageResource(m),
				})
			}
		}

		return problems
	},
}
