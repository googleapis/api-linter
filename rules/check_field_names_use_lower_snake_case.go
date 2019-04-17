package rules

import (
	"github.com/golang/protobuf/v2/reflect/protoreflect"
	"github.com/jgeewax/api-linter/lint"
	"github.com/jgeewax/api-linter/protohelpers"
)

func init() {
	registerRule(
		&protohelpers.DescriptorCallbacks{
			RuleInfo: lint.NewRuleInfo(
				"naming_format::field_names_use_lower_snake_case",
				"check that field names use lower snake case",
				"https://g3doc.corp.google.com/google/api/tools/linter/g3doc/rules/naming-format.md?cl=head",
				[]lint.FileType{lint.ProtoFile},
				lint.CategorySuggestion,
			),
			FieldDescriptorCallback: checkFieldNamesUseLowerSnakeCase,
		},
	)
}

func checkFieldNamesUseLowerSnakeCase(d protoreflect.FieldDescriptor, s lint.DescriptorSource) ([]lint.Problem, error) {
	return checkNameFormat(d), nil
}
