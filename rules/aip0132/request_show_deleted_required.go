package aip0132

import (
	"fmt"
	"strings"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
	"github.com/stoewer/go-strcase"
)

// List requests should contain a show_deleted field if the resource supports
// soft delete.
var requestShowDeletedRequired = &lint.MessageRule{
	Name: lint.NewRuleName(132, "request-show-deleted-required"),
	OnlyIf: func(m protoreflect.MessageDescriptor) bool {
		if !utils.IsListRequestMessage(m) {
			return false
		}
		// Check for soft-delete support by getting the resource name
		// from the corresponding response message.
		plural := strings.TrimPrefix(strings.TrimSuffix(string(m.Name()), "Request"), "List")
		if resp := utils.FindMessage(m.ParentFile(), fmt.Sprintf("List%sResponse", plural)); resp != nil {
			if paged := resp.Fields().ByName(protoreflect.Name(strcase.SnakeCase(plural))); paged != nil && paged.Message() != nil {
				singular := paged.Message().Name()
				return utils.FindMethod(m.ParentFile(), "Undelete"+string(singular)) != nil
			}
		}
		return false
	},
	LintMessage: func(m protoreflect.MessageDescriptor) []lint.Problem {
		if m.Fields().ByName("show_deleted") != nil {
			return nil
		}
		return []lint.Problem{{
			Message:    "List requests for resources supporting soft delete must have a `bool show_deleted` field.",
			Descriptor: m,
		}}
	},
}