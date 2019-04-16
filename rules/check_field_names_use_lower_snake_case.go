package rules

import (
	"github.com/golang/protobuf/v2/reflect/protoreflect"
	"github.com/jgeewax/api-linter/lint"
	"github.com/jgeewax/api-linter/protohelpers"
)

func init() {
	registerRule(
		&protohelpers.DescriptorCallbacks{
			RuleInfo: protohelpers.NewRuleInfo(
				"check_naming_formats.field",
				"check_naming_formats.field",
				"check that field names use lower snake case",
				[]lint.FileType{lint.ProtoFile},
				lint.CategorySuggestion,
			),
			FieldDescriptorCallback: func(d protoreflect.FieldDescriptor, s lint.DescriptorSource) ([]lint.Problem, error) {
				return checkNameFormat(d), nil
			},
		},
	)
}
