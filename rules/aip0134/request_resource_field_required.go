package aip0134

import (
	"fmt"

	"github.com/googleapis/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
)

// The create request message should have resource field.
var requestResourceFieldRequired = &lint.MessageRule{
	Name:   lint.NewRuleName(134, "request-resource-field-required"),
	OnlyIf: isUpdateRequestMessage,
	LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
		resourceMsgName := extractResource(m.GetName())
		for _, fieldDesc := range m.GetFields() {
			msgDesc := fieldDesc.GetMessageType()
			if msgDesc != nil && msgDesc.GetName() == resourceMsgName {
				// found the resource field.
				return nil
			}
		}

		// No resource field.
		return []lint.Problem{{
			Message:    fmt.Sprintf("Message %q has no %q type field", m.GetName(), resourceMsgName),
			Descriptor: m,
		}}
	},
}
