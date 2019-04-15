package lint_test

import (
	"testing"

	"github.com/jgeewax/api-linter/lint"
	"github.com/jgeewax/api-linter/lint/mocks"
)

//go:generate mockery -all

func TestRulesRegister(t *testing.T) {
	r1 := new(mocks.Rule)
	r2 := new(mocks.Rule)

	r1.On("Name").Return("a")
	r2.On("Name").Return("b")

	rules, _ := lint.NewRules()
	err := rules.Register(r1, r2)
	if err != nil {
		t.Errorf("Register: return error %v, but want nil", err)
	}

	numRegistered := len(rules.All())
	if numRegistered != 2 {
		t.Errorf("Register: got %d rules, but want %d", numRegistered, 2)
	}

	r1.AssertCalled(t, "Name")
	r2.AssertCalled(t, "Name")
}

func TestRulesRegister_Duplicate(t *testing.T) {
	r1 := new(mocks.Rule)
	r2 := new(mocks.Rule)

	r1.On("Name").Return("a")
	r2.On("Name").Return("a")

	rules, _ := lint.NewRules()
	err := rules.Register(r1, r2)
	if err != lint.ErrDuplicateName {
		t.Errorf("Register: got %v, but want ErrDuplicateName", err)
	}

	numRegistered := len(rules.All())
	if numRegistered != 1 {
		t.Errorf("Register: got %d rules, but want %d", numRegistered, 1)
	}

	r1.AssertCalled(t, "Name")
	r2.AssertCalled(t, "Name")
}
