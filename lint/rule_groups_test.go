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

func TestAIPInternalGroup(t *testing.T) {
	tests := []struct {
		name  string
		aip   int
		group string
	}{
		{"InInternalGroup", 9001, "internal"},
		{"NotInInternalGroup_Exactly9000", 9000, ""},
		{"InInternalGroup_LargeNumber", 10000, "internal"},
		{"NotInInternalGroup_LessThan9000", 8999, ""},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := aipInternalGroup(test.aip); got != test.group {
				t.Errorf("aipInternalGroup(%d) got %s, but want %s", test.aip, got, test.group)
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
	groupOne := func(aip int) string {
		if aip == 1 {
			return "ONE"
		}
		return ""
	}
	groupTwo := func(aip int) string {
		if aip == 2 {
			return "TWO"
		}
		return ""
	}
	groups := []func(int) string{
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
