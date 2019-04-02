package rules

import (
	"strings"

	"github.com/jgeewax/api-linter/protovisit"

	"github.com/golang/protobuf/v2/reflect/protoreflect"
	"github.com/iancoleman/strcase"
	"github.com/jgeewax/api-linter/lint"
)

func init() {
	registerRule(
		metadata{
			Set:         "core",
			Name:        "check_naming_formats",
			Description: "check naming formats",
			URL:         `https://g3doc.corp.google.com/google/api/tools/linter/g3doc/rules/naming-format.md?cl=head`,
			FileTypes:   []lint.FileType{lint.ProtoFile},
			Category:    lint.CategorySuggestion,
		},
		&checkers{
			DescriptorVisitor: protovisit.SimpleDescriptorVisitor{},
			DescriptorCheck:   checkDescriptor,
		},
	)
}

func checkDescriptor(d protoreflect.Descriptor, ctx lint.Context) []lint.Problem {
	switch d.(type) {
	case protoreflect.EnumDescriptor:
		return checkEnumNamesUseUpperCamelCase(d.(protoreflect.EnumDescriptor), ctx)
	case protoreflect.EnumValueDescriptor:
		return checkEnumValueNamesUseUpperSnakeCase(d.(protoreflect.EnumValueDescriptor), ctx)
	case protoreflect.FieldDescriptor:
		return checkFieldNamesUseLowerSnakeCase(d.(protoreflect.FieldDescriptor), ctx)
	case protoreflect.MessageDescriptor:
		return checkMessageNamesUseUpperCamelCase(d.(protoreflect.MessageDescriptor), ctx)
	case protoreflect.MethodDescriptor:
		return checkMethodNamesUseUpperCamelCase(d.(protoreflect.MethodDescriptor), ctx)
	case protoreflect.ServiceDescriptor:
		return checkServiceNamesUseUpperCamelCase(d.(protoreflect.ServiceDescriptor), ctx)
	}
	return []lint.Problem{}
}

func checkFieldNamesUseLowerSnakeCase(f protoreflect.FieldDescriptor, ctx lint.Context) []lint.Problem {
	return checkNamingFormat("field", f, "lower_snake_case", lowerSnakeCase)
}

func checkEnumNamesUseUpperCamelCase(f protoreflect.EnumDescriptor, ctx lint.Context) []lint.Problem {
	return checkNamingFormat("enum", f, "UpperCamelCase", upperCamelCase)
}

func checkEnumValueNamesUseUpperSnakeCase(f protoreflect.EnumValueDescriptor, ctx lint.Context) []lint.Problem {
	return checkNamingFormat("enum value", f, "UPPER_SNAKE_CASE", upperSnakeCase)
}

func checkMethodNamesUseUpperCamelCase(f protoreflect.MethodDescriptor, ctx lint.Context) []lint.Problem {
	return checkNamingFormat("method", f, "UpperCamelCase", upperCamelCase)
}

func checkMessageNamesUseUpperCamelCase(f protoreflect.MessageDescriptor, ctx lint.Context) []lint.Problem {
	return checkNamingFormat("message", f, "UpperCamelCase", upperCamelCase)
}

func checkServiceNamesUseUpperCamelCase(f protoreflect.ServiceDescriptor, ctx lint.Context) []lint.Problem {
	return checkNamingFormat("service", f, "UpperCamelCase", upperCamelCase)
}

func upperCamelCase(s string) string {
	return strcase.ToCamel(s)
}

func upperSnakeCase(s string) string {
	return strings.ToUpper(strcase.ToSnake(s))
}

func lowerSnakeCase(s string) string {
	return strings.ToLower(strcase.ToSnake(s))
}

func checkNamingFormat(element string, desc protoreflect.Descriptor, example string, transform func(string) string) []lint.Problem {
	name := string(desc.Name())
	suggestion := transform(name)
	if name != suggestion {
		return []lint.Problem{
			{
				Message:    element + " name '" + name + "' should use " + example,
				Suggestion: suggestion,
				Descriptor: desc,
			},
		}
	}
	return []lint.Problem{}
}
