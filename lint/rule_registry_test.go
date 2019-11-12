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
	"testing"
)

func TestRuleRegistryRegister(t *testing.T) {
	tests := []struct {
		name      string
		aip       int
		ruleNames []RuleName
		err       error
	}{
		{
			name:      "Registered_Okay",
			aip:       111,
			ruleNames: []RuleName{NewRuleName(111, "a"), NewRuleName(111, "b")},
			err:       nil,
		},
		{
			name:      "InvalidRuleName",
			aip:       111,
			ruleNames: []RuleName{NewRuleName(111, "")},
			err:       errInvalidRuleName,
		},
		{
			name:      "InvalidRuleGroup",
			aip:       111,
			ruleNames: []RuleName{NewRuleName(100, "a")},
			err:       errInvalidRuleGroup,
		},
		{
			name:      "Duplicated",
			aip:       111,
			ruleNames: []RuleName{NewRuleName(111, "a"), NewRuleName(111, "a")},
			err:       errDuplicatedRuleName,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rules := []ProtoRule{}
			for _, name := range test.ruleNames {
				rules = append(rules, &FileRule{Name: name})
			}

			registry := NewRuleRegistry()
			err := registry.Register(test.aip, rules...)
			if err != test.err {
				t.Errorf("Register(): got %v, but want %v", err, test.err)
			}
		})
	}
}
