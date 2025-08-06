package main

import (
	"testing"

	"github.com/googleapis/api-linter/lint"
	"github.com/jhump/protoreflect/desc/builder"
)

func TestJulioPrefix(t *testing.T) {
	// Create a test registry
	registry := lint.NewRuleRegistry()

	// Add our custom rule
	if err := AddCustomRules(registry); err != nil {
		t.Fatalf("Failed to register custom rules: %v", err)
	}

	// Define test cases
	tests := []struct {
		name          string
		messageName   string
		expectProblem bool
	}{
		{
			name:          "Valid",
			messageName:   "JulioGetBookRequest",
			expectProblem: false,
		},
		{
			name:          "Invalid",
			messageName:   "GetBookRequest",
			expectProblem: true,
		},
	}

	// Get the rule
	julioPrefix := CreateJulioPrefixRule()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create a message descriptor
			msgBuilder := builder.NewMessage(test.messageName)
			msgBuilder.AddField(builder.NewField("name", builder.FieldTypeString()))

			fileBuilder := builder.NewFile("test.proto")
			fileBuilder.AddMessage(msgBuilder)

			fd, err := fileBuilder.Build()
			if err != nil {
				t.Fatalf("Failed to build test file descriptor: %v", err)
			}

			msg := fd.GetMessageTypes()[0]

			// Apply our rule
			problems := julioPrefix.LintMessage(msg)

			if test.expectProblem && len(problems) == 0 {
				t.Errorf("Expected problems but found none for %s", msg.GetName())
			}
			if !test.expectProblem && len(problems) > 0 {
				t.Errorf("Unexpected problems for %s: %v", msg.GetName(), problems)
			}
		})
	}
}
