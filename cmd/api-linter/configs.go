package main

import "github.com/googleapis/api-linter/lint"

// Register default configuration.
func configs() lint.Configs {
	configs := lint.Configs{
		lint.Config{
			IncludedPaths: []string{"**/*.proto"},
			RuleConfigs: map[string]lint.RuleConfig{
				"core": {},
			},
		},
	}
	return configs
}
