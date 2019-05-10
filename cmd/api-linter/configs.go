package main

import "github.com/googleapis/api-linter/lint"

// Register default configuration.
func configs() lint.RuntimeConfigs {
	configs := lint.RuntimeConfigs{
		lint.RuntimeConfig{
			IncludedPaths: []string{"**/*.proto"},
			RuleConfigs: map[string]lint.RuleConfig{
				"core": {
					Status:   lint.Enabled,
					Category: lint.Warning,
				},
			},
		},
	}
	return configs
}
