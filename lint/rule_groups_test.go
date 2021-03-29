package lint

import (
	"testing"
)

func TestAIPCoreGroup(t *testing.T) {
	tests := []struct {
		name  string
		aip   int
		group string
	}{
		{"InCoreGroup", 1, "core"},
		{"NotInCoreGroup_AIP<=0", 0, ""},
		{"NotInCoreGroup_AIP>=1000", 1000, ""},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := aipCoreGroup(test.aip); got != test.group {
				t.Errorf("aipCoreGroup(%d) got %s, but want %s", test.aip, got, test.group)
			}
		})
	}
}

func TestAIPClientLibrariesGroup(t *testing.T) {
	tests := []struct {
		name  string
		aip   int
		group string
	}{
		{"InClientLibrariesGroup", 4232, "client-libraries"},
		{"NotInClientLibrariesGroup_AIP>=4300", 4300, ""},
		{"NotInClientLibrariesGroup_AIP<4200", 4000, ""},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := aipClientLibrariesGroup(test.aip); got != test.group {
				t.Errorf("aipClientLibrariesGroup(%d) got %s, but want %s", test.aip, got, test.group)
			}
		})
	}
}

func TestGetRuleGroupPanic(t *testing.T) {
	var groups []func(int) string
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("getRuleGroup did not panic")
		}
	}()
	getRuleGroup(0, groups)
}

func TestGetRuleGroup(t *testing.T) {
	var groupOne = func(aip int) string {
		if aip == 1 {
			return "ONE"
		}
		return ""
	}
	var groupTwo = func(aip int) string {
		if aip == 2 {
			return "TWO"
		}
		return ""
	}
	var groups = []func(int) string{
		groupOne,
		groupTwo,
	}

	tests := []struct {
		name  string
		aip   int
		group string
	}{
		{"GroupOne", 1, "ONE"},
		{"GroupTwo", 2, "TWO"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := getRuleGroup(test.aip, groups); got != test.group {
				t.Errorf("getRuleGroup(%d) got %s, but want %s", test.aip, got, test.group)
			}
		})
	}
}
