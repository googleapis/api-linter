// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lint

import (
	"fmt"
)

// RuleRegistry is a registry for registering and looking up rules.
type RuleRegistry map[RuleName]Rule

// Copy returns a new copy of the rules.
func (r RuleRegistry) Copy() RuleRegistry {
	n := make(RuleRegistry, len(r))
	for k, v := range r {
		n[k] = v
	}
	return n
}

// Register registers the list of rules.
// Return an error if any of the rules is found duplicate in the registry.
func (r RuleRegistry) Register(rules ...Rule) error {
	for _, rl := range rules {
		if !rl.GetName().IsValid() {
			return fmt.Errorf("%q is not a valid RuleName", rl.GetName())
		}

		if _, found := r[rl.GetName()]; found {
			return fmt.Errorf("duplicate rule name %q", rl.GetName())
		}

		r[rl.GetName()] = rl
	}
	return nil
}

// All returns all rules.
func (r RuleRegistry) All() []Rule {
	rules := make([]Rule, 0, len(r))
	for _, r1 := range r {
		rules = append(rules, r1)
	}
	return rules
}

// NewRuleRegistry returns a rule registry initialized with the given set of rules.
func NewRuleRegistry(rules ...Rule) (RuleRegistry, error) {
	r := make(RuleRegistry, len(rules))
	err := r.Register(rules...)
	return r, err
}
