package aip0134

import (
	"fmt"

	"github.com/googleapis/api-linter/v2/lint"
	"github.com/googleapis/api-linter/v2/rules/internal/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// The create request message should have resource field.
var requestResourceRequired = &lint.MessageRule{
	Name:   lint.NewRuleName(134, "request-resource-required"),
	OnlyIf: utils.IsUpdateRequestMessage,
	LintMessage: func(m protoreflect.MessageDescriptor) []lint.Problem {
		resourceMsgName := extractResource(string(m.Name()))
		for i := 0; i < m.Fields().Len(); i++ {
			fieldDesc := m.Fields().Get(i)
			msgDesc := fieldDesc.Message()
			if msgDesc != nil && string(msgDesc.Name()) == resourceMsgName {
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
