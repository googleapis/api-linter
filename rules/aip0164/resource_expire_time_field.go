package aip0164

import (
	"github.com/googleapis/api-linter/v2/lint"
	"github.com/googleapis/api-linter/v2/rules/internal/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Resources supporting soft delete must have an expire_time field.
var resourceExpireTimeField = &lint.MessageRule{
	Name: lint.NewRuleName(164, "resource-expire-time-field"),
	OnlyIf: func(m protoreflect.MessageDescriptor) bool {
		resource := string(m.Name())
		file, ok := m.Parent().(protoreflect.FileDescriptor)
		if !ok {
			return false
		}
		return utils.FindMethod(file, "Undelete"+resource) != nil
	},
	LintMessage: func(m protoreflect.MessageDescriptor) []lint.Problem {
		// for backwards compatibility, do not lint on expire_time.
		// previously expire_time was the recommended term.
		if m.Fields().ByName("expire_time") != nil {
			return nil
		}
		if m.Fields().ByName("purge_time") != nil {
			return nil
		}
		return []lint.Problem{{
			Message:    "Resources supporting soft delete must have a `google.protobuf.Timestamp purge_time` field.",
			Descriptor: m,
		}}
	},
}
