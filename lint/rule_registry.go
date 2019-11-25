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
	"errors"
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

var errInvalidRuleName = errors.New("not a valid rule name")
var errInvalidRuleGroup = errors.New("invalid rule group")
var errDuplicatedRuleName = errors.New("duplicate rule name")

// Register registers the list of rules of the same AIP.
// Return an error if any of the rules is found duplicate in the registry.
func (r RuleRegistry) Register(aip int, rules ...Rule) error {
	rulePrefix := getRuleGroup(aip) + nameSeparator + fmt.Sprintf("%04d", aip)
	for _, rl := range rules {
		if !rl.Name().IsValid() {
			return errInvalidRuleName
		}

		if !rl.Name().HasPrefix(rulePrefix) {
			return errInvalidRuleGroup
		}

		if _, found := r[rl.Name()]; found {
			return errDuplicatedRuleName
		}

		r[rl.Name()] = rl
	}
	return nil
}

// NewRuleRegistry creates a new rule registry.
func NewRuleRegistry() RuleRegistry {
	return make(RuleRegistry)
}
