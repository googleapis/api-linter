package lint

import (
	"fmt"
)

// Rules is a registry for registering and looking up rules.
type Rules struct {
	ruleMap map[string]Rule
}

// Copy returns a new copy of the rules.
func (r *Rules) Copy() *Rules {
	n := Rules{
		ruleMap: make(map[string]Rule, len(r.ruleMap)),
	}
	for k, v := range r.ruleMap {
		n.ruleMap[k] = v
	}
	return &n
}

// Register registers the list of rules.
// It returns an error if any of the rules is found duplicate
// in the registry.
func (r *Rules) Register(rules ...Rule) error {
	for _, rl := range rules {
		if _, found := r.ruleMap[rl.Info().Name]; found {
			return fmt.Errorf("duplicate rule name `%s`", rl.Info().Name)
		}

		// set default category and status if they are empty.
		if rl.Info().Category == "" {
			rl.Info().Category = DefaultCategory
		}
		if rl.Info().Status == "" {
			rl.Info().Status = DefaultStatus
		}

		r.ruleMap[rl.Info().Name] = rl
	}
	return nil
}

// All returns all rules.
func (r Rules) All() []Rule {
	rules := []Rule{}
	for _, r1 := range r.ruleMap {
		rules = append(rules, r1)
	}
	return rules
}

// NewRules returns a rule registry initialized with the given set of rules.
func NewRules(rules ...Rule) (*Rules, error) {
	r := Rules{
		ruleMap: make(map[string]Rule),
	}
	err := r.Register(rules...)
	return &r, err
}
