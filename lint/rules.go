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

// Rules is a registry for registering and looking up rules.
type Rules struct {
	ruleMap map[RuleID]Rule
}

// Register registers the list of rules.
// It returns `ErrDuplicateID` if any of the rules is found duplicate
// in the registry.
func (r *Rules) Register(rules ...Rule) error {
	for _, rl := range rules {
		if _, found := r.ruleMap[rl.ID()]; found {
			return ErrDuplicateID
		}
		r.ruleMap[rl.ID()] = rl
	}
	return nil
}

// Merge merges another rule registry.
// If any rule is found duplicate, returns `ErrDuplicateID`.
func (r *Rules) Merge(other Rules) error {
	return r.Register(other.All()...)
}

// FindByConfig looks up a list of rules by a RulesConfig
func (r *Rules) FindByConfig(cfg RulesConfig) []Rule {
	rules := []Rule{}
	for _, setConfig := range cfg.RuleSets {
		rules = append(rules, r.findBySet(setConfig)...)
	}
	return rules
}

// All returns all rules.
func (r Rules) All() []Rule {
	rules := []Rule{}
	for _, r1 := range r.ruleMap {
		rules = append(rules, r1)
	}
	return rules
}

func (r *Rules) findBySet(s RuleSetConfig) []Rule {
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
func NewRules(rules ...Rule) (*Rules, error) {
	r := Rules{
		ruleMap: make(map[RuleID]Rule),
	}
	err := r.Register(rules...)
	return &r, err
}
