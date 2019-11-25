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

// Package lint provides lint functions for Google APIs that register rules and user configurations,
// apply those rules to a lint request, and produce lint results.
package lint

import (
	"regexp"
	"sort"
	"strings"
)

// Linter checks API files and returns a list of detected problems.
type Linter struct {
	rules   RuleRegistry
	configs Configs
}

// New creates and returns a linter with the given rules and configs.
func New(rules RuleRegistry, configs Configs) *Linter {
	l := &Linter{
		rules:   rules,
		configs: configs,
	}
	return l
}

// Lint checks a list of descriptors and returns a slice of responses
// with any findings in those descriptors.
func (l *Linter) Lint(descriptors ...Descriptor) []Response {
	results := make(map[string][]Problem)
	for _, d := range descriptors {
		filePath := d.SourceInfo().File().Path()
		for _, rule := range l.rules {
			// Check if the rule is not disabled by comments or configs.
			if ruleIsEnabled(rule, d, aliasMap) && l.configs.IsRuleEnabled(string(rule.Name()), filePath) {
				for _, p := range rule.Lint(d) {
					p.RuleID = rule.Name()
					results[filePath] = append(results[filePath], p)
				}
			}
		}
	}

	filePaths := []string{}
	for path := range results {
		filePaths = append(filePaths, path)
	}
	sort.Strings(filePaths)

	responses := []Response{}
	for _, filePath := range filePaths {
		problems := results[filePath]
		responses = append(responses, Response{
			FilePath: filePath,
			Problems: problems,
		})
	}
	return responses
}

// ruleIsEnabled returns true if the rule is enabled (not disabled by the comments
// for the given descriptor or its file), false otherwise.
func ruleIsEnabled(rule Rule, d Descriptor, aliasMap map[string]string) bool {
	// Some rules have a legacy name. We add it to the check list.
	ruleName := string(rule.Name())
	names := []string{ruleName, aliasMap[ruleName]}

	commentLines := strings.Split(d.SourceInfo().File().Comments(), "\n")
	commentLines = append(commentLines, strings.Split(d.SourceInfo().LeadingComments(), "\n")...)
	disabledRules := []string{}
	for _, commentLine := range commentLines {
		r := extractDisabledRuleName(commentLine)
		if r != "" {
			disabledRules = append(disabledRules, r)
		}
	}

	for _, name := range names {
		if matchRule(name, disabledRules...) {
			return false
		}
	}

	return true
}

var disableRuleNameRegex = regexp.MustCompile(`api-linter:\s*(.+)\s*=\s*disabled`)

func extractDisabledRuleName(commentLine string) string {
	match := disableRuleNameRegex.FindStringSubmatch(commentLine)
	if len(match) > 0 {
		return match[1]
	}
	return ""
}
