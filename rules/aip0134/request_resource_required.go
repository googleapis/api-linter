package aip0134

import (
	"fmt"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// The create request message should have resource field.
var requestResourceRequired = &lint.MessageRule{
	Name:   lint.NewRuleName(134, "request-resource-required"),
	OnlyIf: utils.IsUpdateRequestMessage,
	LintMessage: func(m protoreflect.MessageDescriptor) []lint.Problem {
		resourceMsgName := extractResource(m.Name())
		for _, fieldDesc := range m.Fields() {
			msgDesc := fieldDesc.GetMessageType()
			if msgDesc != nil && msgDesc.Name() == resourceMsgName {
				// found the resource field.
				return nil
			}
		}

		// No resource field.
		return []lint.Problem{{
			Message:    fmt.Sprintf("Message %q has no %q type field", m.Name(), resourceMsgName),
			Descriptor: m,
		}}
	},
}
