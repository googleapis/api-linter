package rules

import (
	"github.com/golang/protobuf/v2/reflect/protoreflect"
	"github.com/jgeewax/api-linter/lint"
	"github.com/jgeewax/api-linter/protohelpers"
)

func init() {
	registerRule(checkNamingFormats())
}

func checkNamingFormats() lint.Rule {
	return &protohelpers.DescriptorCallbacks{
		RuleInfo: lint.RuleInfo{
			Name:        "check_naming_formats.field",
			Description: "check that field names use lower snake case",
			Url:         "https://g3doc.corp.google.com/google/api/tools/linter/g3doc/rules/naming-format.md?cl=head",
			FileTypes:   []lint.FileType{lint.ProtoFile},
			Category:    lint.CategorySuggestion,
		},
		FieldDescriptorCallback: func(d protoreflect.FieldDescriptor, s lint.DescriptorSource) ([]lint.Problem, error) {
			return checkNameFormat(d), nil
		},
	}
}
