package lint

import (
	"fmt"
)

// Rules is a registry for registering and looking up rules.
type Rules map[RuleName]Rule

// Copy returns a new copy of the rules.
func (r Rules) Copy() Rules {
	n := make(Rules, len(r))
	for k, v := range r {
		n[k] = v
	}
	return n
}

// Register registers the list of rules.
// Return an error if any of the rules is found duplicate in the registry.
func (r Rules) Register(rules ...Rule) error {
	for _, rl := range rules {
		if !rl.Info().Name.IsValid() {
			return fmt.Errorf("%q is not a valid RuleName", rl.Info().Name)
		}

		if _, found := r[rl.Info().Name]; found {
			return fmt.Errorf("duplicate rule name %q", rl.Info().Name)
		}

		r[rl.Info().Name] = rl
	}
	return nil
}

// All returns all rules.
func (r Rules) All() []Rule {
	rules := make([]Rule, 0, len(r))
	for _, r1 := range r {
		rules = append(rules, r1)
	}
	return rules
}

// NewRules returns a rule registry initialized with the given set of rules.
func NewRules(rules ...Rule) (Rules, error) {
	r := make(Rules, len(rules))
	err := r.Register(rules...)
	return r, err
}
