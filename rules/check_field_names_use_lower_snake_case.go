package rules

import (
	"fmt"
	"strings"

	pref "github.com/golang/protobuf/v2/reflect/protoreflect"
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/descriptor"
	"github.com/iancoleman/strcase"
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
			FieldCallback: func(d pref.FieldDescriptor, s lint.DescriptorSource) (problems []lint.Problem, err error) {
				fieldName := string(d.Name())
				suggestion := toLowerSnakeCase(fieldName)
				if fieldName != suggestion {
					problems = append(problems, lint.Problem{
						Message:    fmt.Sprintf("field named %q should use lower_snake_case", fieldName),
						Suggestion: suggestion,
						Descriptor: d,
					})
				}
				return
			},
		},
	}
}

// toLowerSnakeCase converts s to lower_snake_case.
func toLowerSnakeCase(s string) string {
	return strings.ToLower(strcase.ToSnake(s))
}
