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

type mockRule struct {
	info       RuleInfo
	lintCalled int
	lintResp   []Problem
	err        error
}

func (r *mockRule) Info() RuleInfo {
	return r.info
}

func (r *mockRule) Lint(Request) ([]Problem, error) {
	r.lintCalled++
	return r.lintResp, r.err
}

func TestRulesRegister(t *testing.T) {
	r1 := &mockRule{info: RuleInfo{Name: "a"}}
	r2 := &mockRule{info: RuleInfo{Name: "b"}}

	rules, _ := NewRules()
	if err := rules.Register(r1, r2); err != nil {
		t.Errorf("Register: return error %v, but want nil", err)
	}

	numRegistered := len(rules.All())
	if numRegistered != 2 {
		t.Errorf("Register: got %d rules, but want %d", numRegistered, 2)
	}
}

func TestRulesRegister_Duplicate(t *testing.T) {
	r1 := &mockRule{info: RuleInfo{Name: "a"}}
	r2 := &mockRule{info: RuleInfo{Name: "a"}}

	rules, _ := NewRules()
	if err := rules.Register(r1, r2); err == nil {
		t.Errorf("Register with duplicate name")
	}

	numRegistered := len(rules.All())
	if numRegistered != 1 {
		t.Errorf("Register: got %d rules, but want %d", numRegistered, 1)
	}
}
