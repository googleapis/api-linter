package rules

import (
	"fmt"
	"strings"

	"github.com/golang/protobuf/v2/reflect/protoreflect"
	"github.com/googleapis/api-linter/lint"
	"github.com/iancoleman/strcase"
)

type nameFormat struct {
	descType, example string
	transform         func(string) string
}

func checkNameFormat(desc protoreflect.Descriptor) []lint.Problem {
	nf := getNameFormat(desc)
	name := string(desc.Name())
	want := nf.transform(name)
	if name != want {
		return []lint.Problem{
			{
				Message:    fmt.Sprintf("%s named %q should use %s", nf.descType, name, nf.example),
				Suggestion: want,
				Descriptor: desc,
			},
		}
	}
	return nil
}

func getNameFormat(desc protoreflect.Descriptor) nameFormat {
	switch desc.(type) {
	case protoreflect.FieldDescriptor:
		return nameFormat{
			descType:  "field",
			example:   "lower_snake_case",
			transform: lowerSnakeCase,
		}
	default:
		return nameFormat{}
	}
}

func lowerSnakeCase(s string) string {
	return strings.ToLower(strcase.ToSnake(s))
}
