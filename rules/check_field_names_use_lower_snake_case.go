package rules

import (
	"github.com/golang/protobuf/v2/reflect/protoreflect"
	"github.com/jgeewax/api-linter/lint"
)

func init() {
	registerProtoRule(
		ruleInfo{
			name:        "check_naming_formats.field",
			description: "check that field names use lower snake case",
			url:         `https://g3doc.corp.google.com/google/api/tools/linter/g3doc/rules/naming-format.md?cl=head`,
			category:    lint.CategorySuggestion,
		},
		protoCheckers{
			CheckFieldDescriptor: func(d protoreflect.FieldDescriptor, s lint.DescriptorSource) ([]lint.Problem, error) {
				return checkNameFormat(d), nil
			},
		},
	)
}
