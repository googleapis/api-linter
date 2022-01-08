package aip0133

import (
	"fmt"
	"strings"

	"github.com/commure/api-linter/lint"
	"github.com/commure/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
	"github.com/stoewer/go-strcase"
)

var requestParentRequired = &lint.MessageRule{
	Name:   lint.NewRuleName(133, "request-parent-required"),
	OnlyIf: isCreateRequestMessage,
	LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
		if m.FindFieldByName("parent") == nil {
			// Sanity check: If the resource has a pattern, and that pattern
			// contains only one variable, then a parent field is not expected.
			//
			// In order to parse out the pattern, we get the resource message
			// from the request, then get the resource annotation from that,
			// and then inspect the pattern there (oy!).
			singular := getResourceMsgNameFromReq(m)
			if field := m.FindFieldByName(strcase.SnakeCase(singular)); field != nil {
				if hasNoParent(field.GetMessageType()) {
					return nil
				}
			}

			// Nope, this is not the unusual case, and a parent field is expected.
			return []lint.Problem{{
				Message:    fmt.Sprintf("Message %q has no `parent` field", m.GetName()),
				Descriptor: m,
			}}
		}

		return nil
	},
}

func hasNoParent(m *desc.MessageDescriptor) bool {
	if resource := utils.GetResource(m); resource != nil {
		for _, pattern := range resource.GetPattern() {
			if strings.Count(pattern, "{") == 1 {
				return true
			}
		}
	}
	return false
}
