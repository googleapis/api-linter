// Package rules contains rules that checks API styles
// in Google Protobuf APIs.
package rules

import (
	"log"

	"github.com/golang/protobuf/v2/reflect/protoreflect"
	"github.com/jgeewax/api-linter/lint"
)

var coreRules, _ = lint.NewRules()

// Rules returns all rules registered in this package.
func Rules() *lint.Rules {
	return coreRules.Copy()
}

type descCheckFunc func(protoreflect.Descriptor, lint.DescriptorSource) ([]lint.Problem, error)

// registerRuleWithDescCheckFunc registers a rule with rule information and
// a descriptor check function for .proto files.
func registerRuleWithDescCheckFunc(ri ruleInfo, c descCheckFunc) {
	if len(ri.fileTypes) == 0 {
		ri.fileTypes = []lint.FileType{lint.ProtoFile}
	}

	r := ruleBase{
		ruleInfo: ri,
		l:        newProtoLinter(ri, c),
	}

	registerRuleBase(r)
}

func registerRuleBase(r ruleBase) {
	if err := coreRules.Register(r); err != nil {
		log.Fatalf("Error when registering rule '%s': %v", r.Name(), err)
	}
}
