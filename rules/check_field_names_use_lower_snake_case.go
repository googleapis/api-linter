package rules

import (
	"github.com/golang/protobuf/v2/reflect/protoreflect"
	"github.com/jgeewax/api-linter/lint"
	"github.com/jgeewax/api-linter/rules/descriptor"
)

func init() {
	registerRules(checkFieldNamesUseLowerSnakeCase())
}

// checkFieldNamesUseLowerSnakeCase returns a lint.Rule
// that checks if a field name is using lower_snake_case.
func checkFieldNamesUseLowerSnakeCase() lint.Rule {
	return &descriptor.CallbackRule{
		RuleInfo: lint.RuleInfo{
			Name:        lint.NewRuleName("naming_formats", "field_names"),
			Description: "check that field names use lower snake case",
			URI:         "https://g3doc.corp.google.com/google/api/tools/linter/g3doc/rules/naming-format.md?cl=head",
			FileTypes:   []lint.FileType{lint.ProtoFile},
		},
		Callback: descriptor.Callbacks{
			FieldCallback: func(d protoreflect.FieldDescriptor, s lint.DescriptorSource) ([]lint.Problem, error) {
				return checkNameFormat(d), nil
			},
		},
	}
}
