package aip0133

import (
	"fmt"

	"github.com/googleapis/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
)

var requestRequiredFields = &lint.MessageRule{
	Name:   lint.NewRuleName(133, "request-required-fields"),
	OnlyIf: isCreateRequestMessage,
	LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
		if m.FindFieldByName("parent") == nil {
			return []lint.Problem{{
				Message:    fmt.Sprintf("Message %q has no `parent` field", m.GetName()),
				Descriptor: m,
			}}
		}

		return nil
	},
}
