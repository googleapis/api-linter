package rules

import (
	"github.com/golang/protobuf/v2/reflect/protoreflect"
	"github.com/jgeewax/api-linter/lint"
)

func init() {
	registerProtoRule(
		ruleInfo{
			name:        "check_proto_syntax_version",
			description: "check that syntax is proto3",
			url:         `https://g3doc.corp.google.com/google/api/tools/linter/g3doc/rules/proto-version.md?cl=head`,
			category:    lint.CategoryError,
		},
		protoCheckers{
			CheckFileDescriptor: func(f protoreflect.FileDescriptor, s lint.DescriptorSource) ([]lint.Problem, error) {
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
				return nil, nil
			},
		},
	)
}
