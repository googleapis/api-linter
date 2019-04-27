package corp

import "github.com/jgeewax/api-linter/lint"

func Rules() (lint.Rules, error) {
	return lint.NewRules([]lint.Rule{
		// Enforce that all proto files include a version number.
		protoFilesMustIncludeVersion(),
		// Add new rules above this line
	}...)
}
