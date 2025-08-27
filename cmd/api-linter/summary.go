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
	"fmt"
	"sort"

	"github.com/googleapis/api-linter/v2/lint"
	"github.com/olekukonko/tablewriter"
)

// printSummaryTable returns a summary table of violation counts.
func printSummaryTable(responses []lint.Response) ([]byte, error) {
	s := createSummary(responses)

	data := []summary{}
	for ruleID, fileViolations := range s {
		totalViolations := 0
		for _, count := range fileViolations {
			totalViolations += count
		}
		data = append(data, summary{ruleID, totalViolations, len(fileViolations)})
	}
	sort.SliceStable(data, func(i, j int) bool { return data[i].violations < data[j].violations })

	var buf bytes.Buffer
	table := tablewriter.NewWriter(&buf)
	table.SetHeader([]string{"Rule", "Total Violations", "Violated Files"})
	table.SetCaption(true, fmt.Sprintf("Linted %d proto files", len(responses)))
	for _, d := range data {
		table.Append([]string{
			d.ruleID,
			fmt.Sprintf("%d", d.violations),
			fmt.Sprintf("%d", d.files),
		})
	}
	table.Render()

	return buf.Bytes(), nil
}

func createSummary(responses []lint.Response) map[string]map[string]int {
	summary := make(map[string]map[string]int)
	for _, r := range responses {
		filePath := string(r.FilePath)
		for _, p := range r.Problems {
			ruleID := string(p.RuleID)
			if summary[ruleID] == nil {
				summary[ruleID] = make(map[string]int)
			}
			summary[ruleID][filePath]++
		}
	}
	return summary
}

type summary struct {
	ruleID     string
	violations int
	files      int
}
