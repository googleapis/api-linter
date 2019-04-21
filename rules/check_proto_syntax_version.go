package rules

import (
	"github.com/golang/protobuf/v2/reflect/protoreflect"
	"github.com/jgeewax/api-linter/lint"
	"github.com/jgeewax/api-linter/rules/descriptor"
)

func init() {
	registerRules(checkProtoVersion())
}

// checkProtoVersion returns a lint.Rule
// that checks if an API is using proto3.
func checkProtoVersion() lint.Rule {
	return &descriptor.CallbackRule{
		RuleInfo: lint.RuleInfo{
			Name:        "proto_version",
			Description: "APIs should use proto3",
			URI:         `https://g3doc.corp.google.com/google/api/tools/linter/g3doc/rules/proto-version.md?cl=head`,
			FileTypes:   []lint.FileType{lint.ProtoFile},
			Category:    lint.CategoryError,
		},
		Callback: descriptor.Callbacks{
			FileCallback: func(f protoreflect.FileDescriptor, s lint.DescriptorSource) ([]lint.Problem, error) {
				location, _ := s.SyntaxLocation()
				if f.Syntax() != protoreflect.Proto3 {
					return []lint.Problem{
						{
							Message:    "APIs should use proto3",
							Suggestion: "proto3",
							Location:   location,
						},
					}, nil
				}
				return nil, nil
			},
		},
	}
}
