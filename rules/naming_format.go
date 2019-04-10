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
			EnumCheck:         checkEnumNamesUseUpperCamelCase,
			EnumValueCheck:    checkEnumValueNamesUseUpperSnakeCase,
			FieldCheck:        checkFieldNamesUseLowerSnakeCase,
			MessageCheck:      checkMessageNamesUseUpperCamelCase,
			MethodCheck:       checkMethodNamesUseUpperCamelCase,
			ServiceCheck:      checkServiceNamesUseUpperCamelCase,
		},
	)
}

func checkFieldNamesUseLowerSnakeCase(f protoreflect.FieldDescriptor, descriptorSource lint.DescriptorSource) []lint.Problem {
	return checkNamingFormat("field", f, "lower_snake_case", lowerSnakeCase)
}

func checkEnumNamesUseUpperCamelCase(f protoreflect.EnumDescriptor, descriptorSource lint.DescriptorSource) []lint.Problem {
	return checkNamingFormat("enum", f, "UpperCamelCase", upperCamelCase)
}

func checkEnumValueNamesUseUpperSnakeCase(f protoreflect.EnumValueDescriptor, descriptorSource lint.DescriptorSource) []lint.Problem {
	return checkNamingFormat("enum value", f, "UPPER_SNAKE_CASE", upperSnakeCase)
}

func checkMethodNamesUseUpperCamelCase(f protoreflect.MethodDescriptor, descriptorSource lint.DescriptorSource) []lint.Problem {
	return checkNamingFormat("method", f, "UpperCamelCase", upperCamelCase)
}

func checkMessageNamesUseUpperCamelCase(f protoreflect.MessageDescriptor, descriptorSource lint.DescriptorSource) []lint.Problem {
	return checkNamingFormat("message", f, "UpperCamelCase", upperCamelCase)
}

func checkServiceNamesUseUpperCamelCase(f protoreflect.ServiceDescriptor, descriptorSource lint.DescriptorSource) []lint.Problem {
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
