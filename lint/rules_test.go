package lint_test

import (
	"sort"
	"testing"

	"github.com/jgeewax/api-linter/lint"
	"github.com/jgeewax/api-linter/lint/mocks"
)

//go:generate mockery -all

func TestRulesRegister(t *testing.T) {
	r1 := new(mocks.Rule)
	r2 := new(mocks.Rule)

	r1.On("ID").Return(lint.RuleID{"a", "a"})
	r2.On("ID").Return(lint.RuleID{"a", "b"})

	rules, _ := lint.NewRules()
	err := rules.Register(r1, r2)
	if err != nil {
		t.Errorf("Register: return error %v, but want nil", err)
	}

	numRegistered := len(rules.All())
	if numRegistered != 2 {
		t.Errorf("Register: got %d rules, but want %d", numRegistered, 2)
	}

	r1.AssertCalled(t, "ID")
	r2.AssertCalled(t, "ID")
}

func TestRulesRegister_Duplicate(t *testing.T) {
	r1 := new(mocks.Rule)
	r2 := new(mocks.Rule)

	r1.On("ID").Return(lint.RuleID{"a", "a"})
	r2.On("ID").Return(lint.RuleID{"a", "a"})

	rules, _ := lint.NewRules()
	err := rules.Register(r1, r2)
	if err != lint.ErrDuplicateID {
		t.Errorf("Register: got %v, but want ErrDuplicateID", err)
	}

	numRegistered := len(rules.All())
	if numRegistered != 1 {
		t.Errorf("Register: got %d rules, but want %d", numRegistered, 1)
	}

	r1.AssertCalled(t, "ID")
	r2.AssertCalled(t, "ID")
}

func TestRulesFindByConfig(t *testing.T) {
	r1 := new(mocks.Rule)
	r2 := new(mocks.Rule)
	r3 := new(mocks.Rule)

	id1 := lint.RuleID{"a", "a"}
	id2 := lint.RuleID{"a", "b"}
	id3 := lint.RuleID{"c", "d"}

	r1.On("ID").Return(id1)
	r2.On("ID").Return(id2)
	r3.On("ID").Return(id3)

	rules, _ := lint.NewRules(r1, r2, r3)

	tests := []struct {
		cfg   lint.RulesConfig
		rules []string
	}{
		{
			cfg: lint.RulesConfig{
				RuleSets: []lint.RuleSetConfig{
					{
						Set:           "a",
						ExcludedRules: []string{"b"},
					},
				},
			},
			rules: []string{"a:a"},
		},
		{
			cfg: lint.RulesConfig{
				RuleSets: []lint.RuleSetConfig{
					{
						Set:           "a",
						ExcludedRules: []string{"a", "b"},
					},
				},
			},
			rules: []string{},
		},
		{
			cfg: lint.RulesConfig{
				RuleSets: []lint.RuleSetConfig{
					{
						Set:           "c",
						ExcludedRules: []string{},
					},
				},
			},
			rules: []string{"c:d"},
		},
		{
			cfg:   lint.RulesConfig{},
			rules: []string{},
		},
		{
			cfg: lint.RulesConfig{
				RuleSets: []lint.RuleSetConfig{
					{
						Set:           "z",
						ExcludedRules: []string{},
					},
				},
			},
			rules: []string{},
		},
	}

	for _, test := range tests {
		founds := rules.FindByConfig(test.cfg)
		ids := []string{}
		for _, r := range founds {
			ids = append(ids, r.ID().Set+":"+r.ID().Name)
		}
		sort.Strings(ids)
		if !stringsEqual(ids, test.rules) {
			t.Errorf("FindByConfig: got %v, wanted %v", ids, test.rules)
		}
	}
}

func stringsEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
