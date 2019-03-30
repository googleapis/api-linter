package lint

import (
	"errors"
)

var (
	// ErrDuplicateID is the returned error when a duplicated rule ID is
	// found in a rule registry.
	ErrDuplicateID = errors.New("rule registry: found duplicate id")
	// ErrNotFound is the returned error when a rule is not found in a
	// rule registry
	ErrNotFound = errors.New("rule registry: rule is not found")
)

// Rules is a registry for looking up or iterating over rules.
type Rules struct {
	ruleMap map[RuleID]Rule
}

// Register registers the list of rules.
// It returns `ErrDuplicateID` if any rule is found duplicate
// in the registry.
func (r *Rules) Register(rules ...Rule) error {
	for _, rl := range rules {
		found := false
		_, found = r.ruleMap[rl.ID()]
		if found {
			return ErrDuplicateID
		}
		r.ruleMap[rl.ID()] = rl
	}
	return nil
}

// Merge merges another rule registry.
// If any rule is found duplicate, returns `ErrDuplicateID`.
func (r *Rules) Merge(other Rules) error {
	return r.Register(other.AllRules()...)
}

// FindRulesByConfig looks up sets of rules.
func (r *Rules) FindRulesByConfig(cfg RulesConfig) []Rule {
	rules := []Rule{}
	for _, setConfig := range cfg.RuleSets {
		rules = append(rules, r.findRulesBySet(setConfig)...)
	}
	return rules
}

// AllRules returns all rules.
func (r Rules) AllRules() []Rule {
	rules := []Rule{}
	for _, r1 := range r.ruleMap {
		rules = append(rules, r1)
	}
	return rules
}

func (r *Rules) findRulesBySet(s RuleSetConfig) []Rule {
	excludedRules := make(map[string]bool)
	for _, ruleName := range s.ExcludedRules {
		excludedRules[ruleName] = true
	}

	rules := []Rule{}
	for _, r1 := range r.ruleMap {
		if r1.ID().Set == s.Set && !excludedRules[r1.ID().Name] {
			rules = append(rules, r1)
		}
	}
	return rules
}

// NewRules returns a rule registry initialized with the given set of rules.
func NewRules(rules ...Rule) *Rules {
	r := Rules{
		ruleMap: make(map[RuleID]Rule),
	}
	r.Register(rules...)
	return &r
}
