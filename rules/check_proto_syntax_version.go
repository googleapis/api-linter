package rules

import (
	"github.com/golang/protobuf/v2/reflect/protoreflect"
	"github.com/jgeewax/api-linter/lint"
)

func init() {
	registerRuleWithDescCheckFunc(
		ruleInfo{
			Name:        "check_proto_syntax_version",
			Description: "check that syntax is proto3",
			URL:         `https://g3doc.corp.google.com/google/api/tools/linter/g3doc/rules/proto-version.md?cl=head`,
			Category:    lint.CategoryError,
		},
		func(d protoreflect.Descriptor, s lint.DescriptorSource) ([]lint.Problem, error) {
			if f, ok := d.(protoreflect.FileDescriptor); ok {
				location, _ := s.SyntaxLocation()
				if f.Syntax() != protoreflect.Proto3 {
					return []lint.Problem{
						{
							Message:    "Google APIs should use proto3",
							Suggestion: "proto3",
							Location:   location,
						},
					}, nil
				}
			}
			return nil, nil
		},
	)
}
