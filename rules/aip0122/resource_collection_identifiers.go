// Copyright 2022 Google LLC
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
	"unicode"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

var resourceCollectionIdentifiers = &lint.MessageRule{
	Name: lint.NewRuleName(122, "resource-collection-identifiers"),
	OnlyIf: func(m *desc.MessageDescriptor) bool {
		return utils.GetResource(m) != nil
	},
	LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
		var problems []lint.Problem
		resource := utils.GetResource(m)
		for _, p := range resource.GetPattern() {
			if strings.IndexRune(p, '/') == 0 {
				return append(problems, lint.Problem{
					Message:    "Resource patterns must not start with a slash.",
					Descriptor: m,
					Location:   locations.MessageResource(m),
				})
			}

			segs := strings.Split(p, "/")
			for _, seg := range segs {
				// Get first rune of each pattern segment.
				c := []rune(seg)[0]

				if unicode.IsLetter(c) && unicode.IsUpper(c) {
					problems = append(problems, lint.Problem{
						Message:    "Resource patterns must use lowerCamelCase for collection identifiers.",
						Descriptor: m,
						Location:   locations.MessageResource(m),
					})
				}
			}
		}

		return problems
	},
}
