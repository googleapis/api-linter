package lint

import "github.com/golang/protobuf/v2/reflect/protoreflect"

// Response describes the result returned by a rule.
type Response struct {
	FilePath string    `json:"file_path" yaml:"file_path"`
	Problems []Problem `json:"problems" yaml:"problems"`
}

// Problem contains information about a result produced by an API Linter.
type Problem struct {
	// Message provides a short description of the problem.
	Message string `json:"message" yaml:"message"`
	// Suggestion provides a suggested fix, if applicable.
	Suggestion string `json:"suggestion,omitempty" yaml:"suggestion,omitempty"`
	// Location provides the location of the problem. If both
	// `Location` and `Descriptor` are specified, the location
	// is then used from `Location` instead of `Descriptor`.
	Location Location `json:"location" yaml:"location"`
	// Descriptor provides the descriptor related
	// to the problem. If present and `Location` is not
	// specified, then the starting location of the descriptor
	// is used as the location of the problem.
	Descriptor protoreflect.Descriptor `json:"-" yaml:"-"`

	// RuleID provides the ID of the rule that this problem belongs to.
	RuleID RuleName `json:"rule_id" yaml:"rule_id"`

	// The following fields will be set by users.
	Category string `json:"category,omitempty" yaml:"category,omitempty"`
}
