package main

import "github.com/googleapis/api-linter/lint"

// Given lint responses, returns
// 1. A summary map where
// 		key = rule_id, value = number of unique files that violated rule
// 2. length of the rule_id of the longest rule added
func createSummary(responses []lint.Response) (map[string]int, int){
	summaryWithFileName, longestRuleLen := createSummaryIncludingFilename(responses)
	summary := make(map[string]int)
	for ruleID, filePaths := range summaryWithFileName {
		summary[ruleID] = len(filePaths)
	}
	return summary, longestRuleLen
}

// Given lint responses, returns
// 1. A summary map where
// 		key = rule_id, value = set of filePaths violating the rule
// 		set of filePaths is a map[string]bool where key=filename
// 		and value is a bool indicating it's present
// 2. length of the rule_id of the longest rule added
func createSummaryIncludingFilename(responses []lint.Response) (map[string]map[string]bool, int) {
	summary := make(map[string]map[string]bool)
	longestRuleLen := 0
	for _, response := range responses {
		pathToAdd := string(response.FilePath)
		problems := response.Problems
		for _, currentProb := range problems {
			ruleName := string(currentProb.RuleID)
			if existingPaths, ok := summary[ruleName]; ok {
				if _, is_exist := existingPaths[pathToAdd]; !is_exist {
					existingPaths[pathToAdd] = true
				}
			} else {
				longestRuleLen = max(longestRuleLen, len(ruleName))
				summary[ruleName] = map[string]bool{pathToAdd: true}
			}

		}
	}
	return summary, longestRuleLen
}

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}