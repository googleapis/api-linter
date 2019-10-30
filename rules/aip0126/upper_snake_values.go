package aip0126

import (
	"fmt"
	"strings"

	"github.com/googleapis/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
	"github.com/stoewer/go-strcase"
)

// All enum values must use UPPER_SNAKE_CASE.
var enumValueUpperSnakeCase = &lint.EnumRule{
	Name: lint.NewRuleName("core", "0126", "upper-snake-values"),
	URI:  "https://aip.dev/126#guidance",
	LintEnum: func(e *desc.EnumDescriptor) []lint.Problem {
		var problems []lint.Problem
		for _, v := range e.GetValues() {
			if got, want := v.GetName(), toUpperSnakeCase(v.GetName()); got != want {
				problems = append(problems, lint.Problem{
					Message:    fmt.Sprintf("Enum value %q must use UPPER_SNAKE_CASE.", got),
					Suggestion: want,
					Descriptor: v,
					Location:   lint.DescriptorNameLocation(v),
				})
			}
		}
		return problems
	},
}

func toUpperSnakeCase(s string) string {
	return strings.ToUpper(strcase.SnakeCase(s))
}
