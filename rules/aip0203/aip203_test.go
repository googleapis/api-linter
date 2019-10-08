package aip0203

import (
	"strings"
	"testing"

	"github.com/googleapis/api-linter/lint"
)

func TestAddRules(t *testing.T) {
	rules := make(lint.RuleRegistry)
	AddRules(rules)
	for ruleName := range rules {
		if !strings.HasPrefix(string(ruleName), "core::0203") {
			t.Errorf("Rule %s is not namespaced to core::0203.", ruleName)
		}
	}
}
