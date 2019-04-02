package lint

// RulesConfig describes which rule sets to use.
type RulesConfig struct {
	RuleSets []RuleSetConfig
}

// RuleSetConfig describes which rules to be excluded in a set.
type RuleSetConfig struct {
	Set           string
	ExcludedRules []string
}
