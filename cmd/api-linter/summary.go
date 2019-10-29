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

package main

import (
	"bytes"
	"strconv"
	"text/template"

	"github.com/googleapis/api-linter/lint"
)

// emitSummary will return a []byte that shows the percentage of files failing for each rule
func emitSummary(responses []lint.Response) ([]byte, error) {
	var buffer bytes.Buffer
	summary := createSummary(responses)
	defaultSummaryColWidth := 25
	colOneFormat := "%-" + strconv.Itoa(max(summary.LongestRuleLen+5, defaultSummaryColWidth)) + "s"
	colTwoFormat := "%" + strconv.Itoa(defaultSummaryColWidth) + "s"

	lintSummaryTemplate, err := template.New("lintSummary").Funcs(
		template.FuncMap{"calcPercentage": func(filePaths map[string]bool, numSourceFiles int) float64 {
			return float64(len(filePaths)) / float64(numSourceFiles) * 100
		}}).Parse(`
----------SUMMARY TABLE---------
{{ printf "Linted %d proto files." .Summary.LongestRuleLen | printf .ColOneFormat }}
{{ printf .ColOneFormat "Rule"}} {{ printf .ColTwoFormat "Violations (Percent)" -}}
{{$colOneFormat := .ColOneFormat -}}
{{$colTwoFormat := .ColTwoFormat -}}
{{$numSourceFiles := .Summary.NumSourceFiles}}
{{range $ruleID, $filePaths := .Summary.Violations -}}
{{printf 	$colOneFormat $ruleID}}{{printf "%d (%.2f%%)" (len $filePaths) (calcPercentage $filePaths $numSourceFiles) | printf $colTwoFormat}}
{{end}}
`)
	if err != nil {
		panic(err)
	}
	templateData := LintSummaryTemplateData{colOneFormat, colTwoFormat, summary}
	err = lintSummaryTemplate.Execute(&buffer, templateData)
	if err != nil {
		panic(err)
	}
	return buffer.Bytes(), nil
}

func createSummary(responses []lint.Response) (summary LintSummary) {
	summary.NumSourceFiles = len(responses)
	summary.Violations = make(map[string]map[string]bool)
	for _, response := range responses {
		pathToAdd := string(response.FilePath)
		problems := response.Problems
		for _, currentProb := range problems {
			ruleName := string(currentProb.RuleID)
			if existingPaths, ok := summary.Violations[ruleName]; ok {
				if _, isExist := existingPaths[pathToAdd]; !isExist {
					existingPaths[pathToAdd] = true
				}
			} else {
				summary.LongestRuleLen = max(summary.LongestRuleLen, len(ruleName))
				summary.Violations[ruleName] = map[string]bool{pathToAdd: true}
			}
		}
	}
	return
}

// LintSummary summarizes a lint run, including which files have violations.
type LintSummary struct {
	// key = rule_id, value = set of unique files that violated rule
	Violations map[string]map[string]bool
	// length of the rule_id of the longest rule added
	LongestRuleLen int
	// count of files from the original source.
	NumSourceFiles int
}

// LintSummaryTemplateData provides formatting data for the file-level report.
type LintSummaryTemplateData struct {
	ColOneFormat string
	ColTwoFormat string
	Summary      LintSummary
}

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}
