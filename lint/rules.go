package lint

import (
	"errors"
)

var (
	// ErrDuplicateName is the returned error when a duplicated rule ID is
	// found in a rule registry.
	ErrDuplicateName = errors.New("rule registry: found duplicate name")
	// ErrNotFound is the returned error when a rule is not found in a
	// rule registry
	ErrNotFound = errors.New("rule registry: rule is not found")
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
// It returns `ErrDuplicateName` if any of the rules is found duplicate
// in the registry.
func (r *Rules) Register(rules ...Rule) error {
	for _, rl := range rules {
		if _, found := r.ruleMap[rl.Info().Name()]; found {
			return ErrDuplicateName
		}
		r.ruleMap[rl.Info().Name()] = rl
	}
	return nil
}

// Merge merges another rule registry.
// If any rule is found duplicate, returns `ErrDuplicateName`.
func (r *Rules) Merge(other Rules) error {
	return r.Register(other.All()...)
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
