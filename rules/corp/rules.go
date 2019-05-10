package corp

import "github.com/googleapis/api-linter/lint"

// Rules resturns a list of registered rules.
func Rules() (lint.Rules, error) {
	return lint.NewRules([]lint.Rule{
		// Enforce that all proto files include a version number.
		protoFilesMustIncludeVersion(),
		// Add new rules above this line
	}...)
}
