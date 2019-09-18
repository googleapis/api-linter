package aip0126

import (
	"strings"
	"testing"

	"github.com/googleapis/api-linter/lint"
)

func TestAddRules(t *testing.T) {
	rules := make(lint.RuleRegistry)
	AddRules(rules)
	for ruleName := range rules {
		if !strings.HasPrefix(string(ruleName), "core::0126") {
			t.Errorf("Rule %s is not namespaced to core::0126.", ruleName)
		}
	}
}
