// Package rules contains rules that checks API styles
// in Google Protobuf APIs.
package rules

import (
	"log"

	"github.com/jgeewax/api-linter/lint"
)

var coreRules, _ = lint.NewRules()

// Rules returns all rules registered in this package.
func Rules() *lint.Rules {
	return coreRules
}

func registerRule(metadata metadata, l linter) {
	r := ruleBase{
		metadata: metadata,
		l:        l,
	}
	registerRuleBase(r)
}

func registerRuleFunc(metadata metadata, l lintFuncType) {
	r := ruleBase{
		metadata: metadata,
		l:        l,
	}
	registerRuleBase(r)
}

func registerRuleBase(r ruleBase) {
	if err := coreRules.Register(r); err != nil {
		log.Fatalf("Error when registering rule '%s': %v", r.ID(), err)
	}
}

type lintFuncType func(lint.Request) (lint.Response, error)

func (f lintFuncType) Lint(rule lint.Rule, req lint.Request) (lint.Response, error) {
	return f(req)
}
