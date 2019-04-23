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
