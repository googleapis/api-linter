package rules

import (
	"github.com/golang/protobuf/v2/reflect/protoreflect"
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/descriptor"
)

func init() {
	registerRules(checkProtoVersion())
}

// checkProtoVersion returns a lint.Rule
// that checks if an API is using proto3.
func checkProtoVersion() lint.Rule {
	return &descriptor.CallbackRule{
		RuleInfo: lint.RuleInfo{
			Name:         lint.NewRuleName("core", "proto_version"),
			Description:  "APIs should use proto3",
			RequestTypes: []lint.RequestType{lint.ProtoRequest},
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
