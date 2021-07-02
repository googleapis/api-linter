package aip0132

import (
	"fmt"
	"strings"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
	"github.com/stoewer/go-strcase"
)

// List requests should contain a show_deleted field if the resource supports
// soft delete.
var requestShowDeletedRequired = &lint.MessageRule{
	Name: lint.NewRuleName(132, "request-show-deleted-required"),
	OnlyIf: func(m *desc.MessageDescriptor) bool {
		if !isListRequestMessage(m) {
			return false
		}
		// Check for soft-delete support by getting the resource name
		// from the corresponding response message.
		plural := strings.TrimPrefix(strings.TrimSuffix(m.GetName(), "Request"), "List")
		if resp := utils.FindMessage(m.GetFile(), fmt.Sprintf("List%sResponse", plural)); resp != nil {
			if paged := resp.FindFieldByName(strcase.SnakeCase(plural)); paged != nil && paged.GetMessageType() != nil {
				singular := paged.GetMessageType().GetName()
				return utils.FindMethod(m.GetFile(), "Undelete"+singular) != nil
			}
		}
		return false
	},
	LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
		if m.FindFieldByName("show_deleted") != nil {
			return nil
		}
		return []lint.Problem{{
			Message:    "List requests for resources supporting soft delete must have a `bool show_deleted` field.",
			Descriptor: m,
		}}
	},
}
