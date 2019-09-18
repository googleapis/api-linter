package aip0126

import "github.com/googleapis/api-linter/lint"

// AddRules adds all of the AIP-126 rules to the provided registry.
func AddRules(r lint.RuleRegistry) {
	r.Register(
		enumValueUpperSnakeCase,
	)
}
