package main

import (
	"github.com/googleapis/api-linter/lint"
	core "github.com/googleapis/api-linter/rules"
)

// Register rules.
func rules() []lint.Rule {
	var rules []lint.Rule
	rules = append(rules, core.Rules().All()...)
	return rules
}
