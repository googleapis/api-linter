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
	r1 := &FileRule{Name: "a"}
	r2 := &FileRule{Name: "b"}

	rules, _ := NewRuleRegistry()
	if err := rules.Register(r1, r2); err != nil {
		t.Errorf("Register: return error %v, but want nil", err)
	}

	numRegistered := len(rules.All())
	if numRegistered != 2 {
		t.Errorf("Register: got %d rules, but want %d", numRegistered, 2)
	}
}

func TestRulesRegister_Duplicate(t *testing.T) {
	r1 := &FileRule{Name: "a"}
	r2 := &FileRule{Name: "a"}

	rules, _ := NewRuleRegistry()
	if err := rules.Register(r1, r2); err == nil {
		t.Errorf("Register with duplicate name")
	}

	numRegistered := len(rules.All())
	if numRegistered != 1 {
		t.Errorf("Register: got %d rules, but want %d", numRegistered, 1)
	}
}
