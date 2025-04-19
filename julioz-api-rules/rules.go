package main

import (
	"strings"

	"github.com/googleapis/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
)

// CreateJulioPrefixRule creates the rule that enforces message names to start with "Julio"
func CreateJulioPrefixRule() *lint.MessageRule {
	return &lint.MessageRule{
		Name: lint.NewRuleName(9001, "julio-prefix"),
		LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
			name := m.GetName()
			if !strings.HasPrefix(name, "Julio") {
				return []lint.Problem{{
					Message:    "Message names must start with 'Julio'",
					Descriptor: m,
					Suggestion: "Julio" + name,
				}}
			}
			return nil
		},
	}
}
