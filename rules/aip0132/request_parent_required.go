package aip0132

import (
	"fmt"
	"strings"

	"github.com/googleapis/api-linter/v2/lint"
	"github.com/googleapis/api-linter/v2/rules/internal/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
	"github.com/stoewer/go-strcase"
)

// The List standard method should contain a parent field.
var requestParentRequired = &lint.MessageRule{
	Name:   lint.NewRuleName(132, "request-parent-required"),
	OnlyIf: utils.IsListRequestMessage,
	LintMessage: func(m protoreflect.MessageDescriptor) []lint.Problem {
		// Rule check: Establish that a `parent` field is present.
		if m.Fields().ByName("parent") == nil {
			// In order to parse out the pattern, we get the resource message
			// from the response, then get the resource annotation from that,
			// and then inspect the pattern there (oy!).
			plural := strings.TrimPrefix(strings.TrimSuffix(string(m.Name()), "Request"), "List")
			if resp := utils.FindMessage(m.ParentFile(), fmt.Sprintf("List%sResponse", plural)); resp != nil {
				if resField := resp.Fields().ByName(protoreflect.Name(strcase.SnakeCase(plural))); resField != nil {
					if !utils.HasParent(utils.GetResource(resField.Message())) {
						return nil
					}
				}
			}

			// Nope, there should be a parent field and is not. Complain.
			return []lint.Problem{{
				Message:    fmt.Sprintf("Message %q has no `parent` field", m.Name()),
				Descriptor: m,
			}}
		}

		return nil
	},
}