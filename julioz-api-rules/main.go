package main

import (
	"github.com/googleapis/api-linter/lint"
)

// AddCustomRules adds custom rules to the provided registry.
// This is the required entry point for api-linter plugins.
func AddCustomRules(registry lint.RuleRegistry) error {
	// Create the custom rule
	julioPrefix := CreateJulioPrefixRule()

	// Register the rule
	return registry.Register(9001, julioPrefix)
}
