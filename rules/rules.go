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
	return coreRules
}

// DescriptorChecker defines an operation that checks a Descriptor
// and returns a list of Problem and if applicable, an error.
type DescriptorChecker interface {
	Check(protoreflect.Descriptor) ([]lint.Problem, error)
}

// RegisterRuleWithChecker registers a rule with rule information and
// a descriptor checker for .proto files.
func RegisterRuleWithChecker(i RuleInfo, c DescriptorChecker) {
	r := ruleBase{
		RuleInfo: i,
		l:        newProtoLinter(i, c),
	}
	registerRuleBase(r)
}

func registerRuleBase(r ruleBase) {
	if err := coreRules.Register(r); err != nil {
		log.Fatalf("Error when registering rule '%s': %v", r.Name(), err)
	}
}
