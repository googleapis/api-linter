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

import "testing"

func TestRuleName_IsValid(t *testing.T) {
	tests := []struct {
		name  RuleName
		valid bool
	}{
		{"", false},
		{"a:::b", false},
		{"a::::b", false},
		{"a", true},
		{"::my_rule", false},
		{"my_namespace::", false},
		{"my_namespace:my_rule", false},
		{"my_namespace::my_rule", true},
		{"my_rule", true},
	}

	for _, test := range tests {
		if test.name.IsValid() != test.valid {
			t.Errorf("%q.IsValid()=%t; want %t", test.name, test.name.IsValid(), test.valid)
		}
	}
}

func TestNewRuleName(t *testing.T) {
	tests := []struct {
		segments []string
		name     RuleName
	}{
		{[]string{}, ""},
		{[]string{""}, ""},
		{[]string{"my_namespace", "my_rule"}, "my_namespace::my_rule"},
		{[]string{"my", "name", "space", "foo"}, "my::name::space::foo"},
	}

	for _, test := range tests {
		if NewRuleName(test.segments...) != test.name {
			t.Errorf("NewRuleName(%v)=%q; want %q", test.segments, NewRuleName(test.segments...), test.name)
		}
	}
}

func TestRuleName_HasPrefix(t *testing.T) {
	tests := []struct {
		r         RuleName
		prefix    []string
		hasPrefix bool
	}{
		{"a::b::c", []string{"a", "b"}, true},
		{"a::b::c", []string{"a"}, true},
		{"a::b::c", []string{"a::b"}, true},
		{"a::b::c::d", []string{"a::b", "c"}, true},
		{"a::b::c", []string{"a::b::c"}, true},
		{"ab::b::c", []string{"a"}, false},
	}

	for _, test := range tests {
		if test.r.HasPrefix(test.prefix...) != test.hasPrefix {
			t.Errorf(
				"%q.HasPrefix(%v)=%t; want %t",
				test.r, test.prefix, test.r.HasPrefix(test.prefix...), test.hasPrefix,
			)
		}
	}
}

func TestRuleName_parent(t *testing.T) {
	tests := []struct {
		r RuleName
		p RuleName
	}{
		{"a::b::c", "a::b"},
		{"a", ""},
		{"foo::bar::baz::qux", "foo::bar::baz"},
	}

	for _, test := range tests {
		if test.r.parent() != test.p {
			t.Errorf("%q.parent()=%q; want %q", test.r, test.r.parent(), test.p)
		}
	}
}
