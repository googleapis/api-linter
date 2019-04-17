package rules

import (
	"github.com/golang/protobuf/v2/reflect/protoreflect"
	"github.com/jgeewax/api-linter/lint"
	"github.com/jgeewax/api-linter/protohelpers"
)

func init() {
	registerRule(
		&protohelpers.DescriptorCallbacks{
			RuleInfo: lint.NewRuleInfo(
				"proto_syntax_version",
				"check that syntax is proto3",
				`https://g3doc.corp.google.com/google/api/tools/linter/g3doc/rules/proto-version.md?cl=head`,
				[]lint.FileType{lint.ProtoFile},
				lint.CategoryError,
			),
			FileDescriptorCallback: checkProtoSyntaxVersion,
		},
	)
}

func checkProtoSyntaxVersion(f protoreflect.FileDescriptor, s lint.DescriptorSource) ([]lint.Problem, error) {
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
}
