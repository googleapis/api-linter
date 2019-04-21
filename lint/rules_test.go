package lint

import (
	"testing"
)

type mockRule struct {
	info       *RuleInfo
	lintCalled int
	lintResp   Response
	err        error
}

func (r *mockRule) Info() *RuleInfo {
	return r.info
}

func (r *mockRule) Lint(Request) (Response, error) {
	r.lintCalled++
	return r.lintResp, r.err
}

func TestRulesRegister(t *testing.T) {
	r1 := &mockRule{info: &RuleInfo{Name: "a"}}
	r2 := &mockRule{info: &RuleInfo{Name: "b"}}

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
	r1 := &mockRule{info: &RuleInfo{Name: "a"}}
	r2 := &mockRule{info: &RuleInfo{Name: "a"}}

	rules, _ := NewRules()
	if err := rules.Register(r1, r2); err == nil {
		t.Errorf("Register with duplicate name")
	}

	numRegistered := len(rules.All())
	if numRegistered != 1 {
		t.Errorf("Register: got %d rules, but want %d", numRegistered, 1)
	}
}

func TestRulesRegister_Default(t *testing.T) {
	r1 := &mockRule{info: &RuleInfo{Name: "a"}}

	rules, _ := NewRules()
	if err := rules.Register(r1); err != nil {
		t.Errorf("Register: return error %v, but want nil", err)
	}

	if got, want := r1.Info().Category, DefaultCategory; got != want {
		t.Errorf("Register got default category `%s`, but want `%s`", got, want)
	}
	if got, want := r1.Info().Status, DefaultStatus; got != want {
		t.Errorf("Register got default status `%s`, but want `%s`", got, want)
	}
}
