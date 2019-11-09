package rules

import (
	"testing"

	"github.com/googleapis/api-linter/lint"
)

func TestAdd(t *testing.T) {
	if err := Add(lint.NewRuleRegistry()); err != nil {
		t.Errorf("Add got an error: %v", err)
	}
}
