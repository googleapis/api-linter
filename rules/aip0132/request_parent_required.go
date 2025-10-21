package aip0132

import (
	"fmt"
	"strings"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
	"github.com/stoewer/go-strcase"
)

// The List standard method should contain a parent field.
var requestParentRequired = &lint.MessageRule{
	Name:   lint.NewRuleName(132, "request-parent-required"),
	OnlyIf: utils.IsListRequestMessage,
	LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
		// Rule check: Establish that a `parent` field is present.
		if m.FindFieldByName("parent") == nil {
			// In order to parse out the pattern, we get the resource message
			// from the response, then get the resource annotation from that,
			// and then inspect the pattern there (oy!).
			plural := strings.TrimPrefix(strings.TrimSuffix(m.GetName(), "Request"), "List")
			if resp := utils.FindMessage(m.GetFile(), fmt.Sprintf("List%sResponse", plural)); resp != nil {
				if resField := resp.FindFieldByName(strcase.SnakeCase(plural)); resField != nil {
					if !utils.HasParent(utils.GetResource(resField.GetMessageType())) {
						return nil
					}
				}
			}

			// Nope, there should be a parent field and is not. Complain.
			return []lint.Problem{{
				Message:    fmt.Sprintf("Message %q has no `parent` field", m.GetName()),
				Descriptor: m,
			}}
		}

		return nil
	},
}
