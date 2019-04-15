package rules

import (
	"strings"

	"github.com/jgeewax/api-linter/lint"

	"github.com/golang/protobuf/v2/reflect/protoreflect"
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
				Message:    nf.descType + " name '" + name + "' should use " + nf.example,
				Suggestion: want,
				Descriptor: desc,
			},
		}
	}
	return []lint.Problem{}
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
