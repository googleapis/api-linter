package main

import (
	"log"
	"os"

	"github.com/googleapis/api-linter/lint"
	core "github.com/googleapis/api-linter/rules"
)

func main() {
	if err := run(getRules(), getConfigs(), os.Args); err != nil {
		log.Fatal(err)
	}
}

// Register default configuration.
func getConfigs() lint.RuntimeConfigs {
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

// Register rules.
func getRules() []lint.Rule {
	var rules []lint.Rule
	rules = append(rules, core.Rules().All()...)
	return rules
}
