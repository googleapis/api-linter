package aip0126

import (
	"fmt"
	"strings"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/stoewer/go-strcase"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// All enum values must use UPPER_SNAKE_CASE.
var enumValueUpperSnakeCase = &lint.EnumRule{
	Name: lint.NewRuleName(126, "upper-snake-values"),
	LintEnum: func(e protoreflect.EnumDescriptor) []lint.Problem {
		var problems []lint.Problem
		for i := 0; i < e.Values().Len(); i++ {
			v := e.Values().Get(i)
			if got, want := string(v.Name()), toUpperSnakeCase(string(v.Name())); got != want {
				problems = append(problems, lint.Problem{
					Message:    fmt.Sprintf("Enum value %q must use UPPER_SNAKE_CASE.", got),
					Suggestion: want,
					Descriptor: v,
					Location:   locations.DescriptorName(v),
				})
			}
		}
		return problems
	},
}

func toUpperSnakeCase(s string) string {
	return strings.ToUpper(strcase.SnakeCase(s))
}
