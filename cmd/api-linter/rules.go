package main

import (
	"github.com/googleapis/api-linter/lint"
	core "github.com/googleapis/api-linter/rules"
)

// Register rules.
func rules() lint.Rules {
	rules := core.Rules().Copy()
	return rules
}
