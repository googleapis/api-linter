package rules

import (
	"github.com/golang/protobuf/v2/reflect/protoreflect"
	"github.com/jgeewax/api-linter/lint"
)

func init() {
	registerRuleFunc(
		metadata{
			Set:         "core",
			Name:        "check_proto_syntax",
			Description: `Use "proto3" instead of "proto2"`,
			URL:         `https://g3doc.corp.google.com/google/api/tools/linter/g3doc/rules/proto-version.md?cl=head`,
			FileTypes:   []lint.FileType{lint.ProtoFile},
			Category:    lint.CategoryError,
		},
		checkProtoSyntax,
	)
}

func checkProtoSyntax(req lint.Request) (lint.Response, error) {
	f := req.ProtoFile()
	if f.Syntax() != protoreflect.Proto3 {
		return lint.Response{
			Problems: []lint.Problem{
				{
					Message:    "Uses proto3",
					Suggestion: "proto3",
					Location:   findSyntaxLocation(req.DescriptorSource()),
				},
			},
		}, nil
	}
	return lint.Response{}, nil
}

const syntaxTag = 12

func findSyntaxLocation(source lint.DescriptorSource) lint.Location {
	if loc, err := source.SyntaxLocation(); err == nil {
		return loc
	}
	return lint.StartLocation
}
