package lint

import "errors"

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
	ruleMap map[ID]Rule
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

// FindRuleByID looks up any rule by its `ID`.
// If not found, return (nil, ErrNotFound).
func (r *Rules) FindRuleByID(id ID) (Rule, error) {
	rl, found := r.ruleMap[id]
	if !found {
		return nil, ErrNotFound
	}
	return rl, nil
}

// FindRulesBySet looks up a set of rules.
func (r *Rules) FindRulesBySet(set string) []Rule {
	res := []Rule{}
	for _, rl := range r.ruleMap {
		if rl.ID().Set == set {
			res = append(res, rl)
		}
	}
	return res
}

// NewRules returns a rule registry initialized with the given set of rules.
func NewRules(rules ...Rule) *Rules {
	r := Rules{
		ruleMap: make(map[ID]Rule),
	}
	r.Register(rules...)
	return &r
}
