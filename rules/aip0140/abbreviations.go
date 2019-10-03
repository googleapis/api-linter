// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aip0140

import (
	"fmt"
	"strings"

	"github.com/googleapis/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
)

// checkAbbreviations checks for each long name in the descriptor's name,
// and returns problems if one is found that suggests the abbreviation instead.
func checkAbbreviations(d desc.Descriptor, x func(string) string) (problems []lint.Problem) {
	abbv := map[string]string{
		"configuration": "config",
		"identifier":    "id",
		"information":   "info",
		"specification": "spec",
		"statistics":    "stats",
	}

	// Iterate over each abbreviation and determine whether the descriptor's
	// name includes the long name.
	for long, short := range abbv {
		if strings.Contains(d.GetName(), x(long)) {
			problems = append(problems, lint.Problem{
				Message: fmt.Sprintf(
					"Use the common abbreviation %q instead of %q.",
					x(short),
					x(long),
				),
				Suggestion: strings.ReplaceAll(d.GetName(), x(long), x(short)),
				Descriptor: d,
				Location:   lint.DescriptorNameLocation(d),
			})
		}
	}
	return
}

var abbreviationsField = &lint.FieldRule{
	Name: lint.NewRuleName("core", "0140", "field-names", "abbreviations"),
	URI:  "https://aip.dev/140#abbreviations",
	LintField: func(f *desc.FieldDescriptor) []lint.Problem {
		return checkAbbreviations(f, strings.ToLower)
	},
}

var abbreviationsMessage = &lint.MessageRule{
	Name: lint.NewRuleName("core", "0140", "message-names", "abbreviations"),
	URI:  "https://aip.dev/140#abbreviations",
	LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
		return checkAbbreviations(m, strings.Title)
	},
}

var abbreviationsEnum = &lint.EnumRule{
	Name: lint.NewRuleName("core", "0140", "enum-names", "abbreviations"),
	URI:  "https://aip.dev/140#abbreviations",
	LintEnum: func(e *desc.EnumDescriptor) (problems []lint.Problem) {
		problems = append(problems, checkAbbreviations(e, strings.Title)...)
		for _, ev := range e.GetValues() {
			problems = append(problems, checkAbbreviations(ev, strings.ToUpper)...)
		}
		return
	},
}

var abbreviationsService = &lint.ServiceRule{
	Name: lint.NewRuleName("core", "0140", "serivice-names", "abbreviations"),
	URI:  "https://aip.dev/140#abbreviations",
	LintService: func(s *desc.ServiceDescriptor) []lint.Problem {
		return checkAbbreviations(s, strings.Title)
	},
}

var abbreviationsMethod = &lint.MethodRule{
	Name: lint.NewRuleName("core", "0140", "method-names", "abbreviations"),
	URI:  "https://aip.dev/140#abbreviations",
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		return checkAbbreviations(m, strings.Title)
	},
}
