package rules

import (
	"github.com/golang/protobuf/v2/reflect/protoreflect"
	"github.com/jgeewax/api-linter/lint"
)

func init() {
	registerRuleWithDescCheckFunc(
		ruleInfo{
			name:        "check_naming_formats.field",
			description: "check that field names use lower snake case",
			url:         `https://g3doc.corp.google.com/google/api/tools/linter/g3doc/rules/naming-format.md?cl=head`,
			category:    lint.CategorySuggestion,
		},
		func(d protoreflect.Descriptor, s lint.DescriptorSource) ([]lint.Problem, error) {
			if _, ok := d.(protoreflect.FieldDescriptor); ok {
				return checkNameFormat(d), nil
			}
			return nil, nil
		},
	)
}
