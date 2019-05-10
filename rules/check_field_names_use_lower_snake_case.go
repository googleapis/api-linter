package rules

import (
	"github.com/golang/protobuf/v2/reflect/protoreflect"
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/descriptor"
)

func init() {
	registerRules(checkFieldNamesUseLowerSnakeCase())
}

// checkFieldNamesUseLowerSnakeCase returns a lint.Rule
// that checks if a field name is using lower_snake_case.
func checkFieldNamesUseLowerSnakeCase() lint.Rule {
	return &descriptor.CallbackRule{
		RuleInfo: lint.RuleInfo{
			Name:         lint.NewRuleName("core", "naming_formats", "field_names"),
			Description:  "check that field names use lower snake case",
			RequestTypes: []lint.RequestType{lint.ProtoRequest},
		},
		Callback: descriptor.Callbacks{
			FieldCallback: func(d protoreflect.FieldDescriptor, s lint.DescriptorSource) ([]lint.Problem, error) {
				return checkNameFormat(d), nil
			},
		},
	}
}
