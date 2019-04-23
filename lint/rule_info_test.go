package lint

import "testing"

func TestRuleName_IsValid(t *testing.T) {
	tests := []struct {
		name  RuleName
		valid bool
	}{
		{"", false},
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
