package aip0132

import (
	"fmt"

	"github.com/googleapis/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
)

// The List standard method should contain a parent field.
var requestParentRequired = &lint.MessageRule{
	Name:   lint.NewRuleName(132, "request-parent-required"),
	OnlyIf: isListRequestMessage,
	LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
		// Rule check: Establish that a `parent` field is present.
		if m.FindFieldByName("parent") == nil {
			return []lint.Problem{{
				Message:    fmt.Sprintf("Message %q has no `parent` field", m.GetName()),
				Descriptor: m,
			}}
		}

		return nil
	},
}
