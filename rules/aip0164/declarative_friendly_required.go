package aip0164

import (
	"fmt"

	"github.com/commure/api-linter/lint"
	"github.com/commure/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

// Declarative-friendly resources should have an Undelete method.
var declarativeFriendlyRequired = &lint.MessageRule{
	Name: lint.NewRuleName(164, "declarative-friendly-required"),
	OnlyIf: func(m *desc.MessageDescriptor) bool {
		if resource := utils.DeclarativeFriendlyResource(m); resource == m {
			return true
		}
		return false
	},
	LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
		resource := m.GetName()
		want := fmt.Sprintf("Undelete%s", resource)
		delete := fmt.Sprintf("Delete%s", resource)
		if utils.FindMethod(m.GetFile(), want) == nil && utils.FindMethod(m.GetFile(), delete) != nil {
			return []lint.Problem{{
				Message:    fmt.Sprintf("Declarative-friendly %s should have an Undelete method.", resource),
				Descriptor: m,
			}}
		}
		return nil
	},
}
