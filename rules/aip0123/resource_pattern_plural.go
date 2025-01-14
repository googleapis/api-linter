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
)

var resourcePatternPlural = &lint.MessageRule{
	Name: lint.NewRuleName(123, "resource-pattern-plural"),
	OnlyIf: func(m *desc.MessageDescriptor) bool {
		return utils.IsResource(m) && len(utils.GetResource(m).GetPattern()) > 0 && utils.GetResourcePlural(utils.GetResource(m)) != "" && !utils.IsSingletonResource(m)
	},
	LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
		var problems []lint.Problem
		res := utils.GetResource(m)
		nested := isNestedName(res)
		var nn string
		if nested {
			nn = nestedPlural(res)
		}

		patterns := res.GetPattern()
		plural := fmt.Sprintf("/%s/", utils.GetResourcePlural(res))
		if isRootLevelResource(res) {
			plural = strings.TrimPrefix(plural, "/")
		}
		nn = fmt.Sprintf("/%s/", nn)

		// If the first pattern is reduced or non-compliant, but is nested name eligible, we want to recommend the nested name.
		nestedFirstPattern := nested && (strings.Contains(patterns[0], nn) || !strings.Contains(patterns[0], plural))

		for ndx, pattern := range patterns {
			if !strings.Contains(pattern, plural) || ndx > 0 && nestedFirstPattern && !strings.Contains(pattern, nn) {
				// allow the reduced, nested name instead if present
				if nested && strings.Contains(pattern, nn) {
					continue
				}

				want := plural
				if nestedFirstPattern {
					want = nn
				}

				problems = append(problems, lint.Problem{
					Message:    fmt.Sprintf("resource pattern %q collection segment must be the resource plural %q", pattern, want),
					Descriptor: m,
					Location:   locations.MessageResource(m),
				})
			}
		}

		return problems
	},
}
