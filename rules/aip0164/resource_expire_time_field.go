package aip0164

import (
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

// Resources supporting soft delete must have an expire_time field.
var resourceExpireTimeField = &lint.MessageRule{
	Name: lint.NewRuleName(164, "resource-expire-time-field"),
	OnlyIf: func(m *desc.MessageDescriptor) bool {
		resource := m.GetName()
		return utils.FindMethod(m.GetFile(), "Undelete"+resource) != nil
	},
	LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
		// for backwards compatibility, do not lint on expire_time.
		// previously expire_time was the recommended term.
		if m.FindFieldByName("expire_time") != nil {
			return nil
		}
		if m.FindFieldByName("purge_time") != nil {
			return nil
		}
		return []lint.Problem{{
			Message:    "Resources supporting soft delete must have a `google.protobuf.Timestamp purge_time` field.",
			Descriptor: m,
		}}
	},
}
