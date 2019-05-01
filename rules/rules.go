// Package rules contains rules that checks API styles
// in Google Protobuf APIs.
package rules

import (
	"log"

	"github.com/googleapis/api-linter/lint"
)

var coreRules, _ = lint.NewRules()

// Rules returns all rules registered in this package.
func Rules() lint.Rules {
	return coreRules.Copy()
}

func registerRules(r ...lint.Rule) {
	for _, rl := range r {
		if err := coreRules.Register(rl); err != nil {
			log.Fatalf("Error when registering rule '%s': %v", rl.Info().Name, err)
		}
	}
}
