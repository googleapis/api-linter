package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/googleapis/api-linter/lint"
)

func createSummary(responses []lint.Response) (summary LintSummary) {
	summary.numSourceFiles = len(responses)
	summary.violationData = make(map[string]map[string]bool)
	for _, response := range responses {
		pathToAdd := string(response.FilePath)
		problems := response.Problems
		for _, currentProb := range problems {
			ruleName := string(currentProb.RuleID)
			if existingPaths, ok := summary.violationData[ruleName]; ok {
				if _, isExist := existingPaths[pathToAdd]; !isExist {
					existingPaths[pathToAdd] = true
				}
			} else {
				summary.longestRuleLen = max(summary.longestRuleLen, len(ruleName))
				summary.violationData[ruleName] = map[string]bool{pathToAdd: true}
			}
		}
	}
	return summary
}

// Given a pointer to a summary of lint responses, and output location,
// this functions will return a []byte which shows the percentage of files failing for each rule
func emitSummary(summary *LintSummary) ([]byte, error) {
	var sb strings.Builder
	DEFAULT_SUMMARY_COL_WIDTH := 25
	colOneFormat := "%-"+strconv.Itoa(max(summary.longestRuleLen+5, DEFAULT_SUMMARY_COL_WIDTH))+"s"
	colTwoFormat := "%"+strconv.Itoa(DEFAULT_SUMMARY_COL_WIDTH)+"s"
	sb.WriteString("\n----------SUMMARY TABLE---------\n")
	sb.WriteString(fmt.Sprintf(colOneFormat,
		fmt.Sprintf("Linted %d proto files.", summary.numSourceFiles)) + "\n")
	sb.WriteString(fmt.Sprintf(colOneFormat, "Rule") +
			fmt.Sprintf(colTwoFormat, "Violations (Percent)") + "\n")

	for rule_id, filePaths := range summary.violationData {
		sb.WriteString(fmt.Sprintf(colOneFormat, string(rule_id)) +
				fmt.Sprintf(colTwoFormat,
					fmt.Sprintf("%d (%.2f%%)",
						len(filePaths), float64(len(filePaths))/float64(summary.numSourceFiles) * 100,
					),
				) + "\n",
		)
	}
	return []byte(sb.String()), nil
}

type LintSummary struct {
	// key = rule_id, value = number of unique files that violated rule
	violationData map[string]map[string]bool
	// length of the rule_id of the longest rule added
	longestRuleLen int
	// count of files from the original source.
	numSourceFiles int
}

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}
