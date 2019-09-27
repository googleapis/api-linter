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

func TestRuleNameValid(t *testing.T) {
	tests := []struct {
		testName string
		ruleName RuleName
	}{
		{"Lower", "aip"},
		{"LowerNumber", "aip0121"},
		{"LowerNumberKebab", "aip-0121"},
		{"Namespaced", "aip::0121"},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			if !test.ruleName.IsValid() {
				t.Errorf("Rule name %q is invalid; want valid.", test.ruleName)
			}
		})
	}
}

func TestRuleNameInvalid(t *testing.T) {
	tests := []struct {
		testName string
		ruleName RuleName
	}{
		{"EmptyString", ""},
		{"TripleColon", "a:::b"},
		{"QuadrupleColon", "a::::b"},
		{"CapitalLetter", "A"},
		{"LeadingDoubleColon", "::my-rule"},
		{"TrailingDoubleColon", "my-namespace::"},
		{"LeadingHyphen", "-core::aip-0131"},
		{"LeadingSegmentHyphen", "core::-aip-0131"},
		{"OnlyHyphen", "-"},
		{"SingleColon", "core:aip-0131"},
		{"Underscore", "core::aip_0131"},
		{"CamelCase", "myRule"},
		{"PascalCase", "MyRule"},
	}

	for _, test := range tests {
		if test.ruleName.IsValid() {
			t.Errorf("Rule name %q is valid; want invalid.", test.ruleName)
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
		{[]string{"my-namespace", "my-rule"}, "my-namespace::my-rule"},
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
