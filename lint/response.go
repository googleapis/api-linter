package lint

import "github.com/golang/protobuf/v2/reflect/protoreflect"

// Response describes the result returned by a rule.
type Response struct {
	Problems []Problem
}

// Problem contains information about a result produced by an API Linter.
type Problem struct {
	// Message provides a short description of the problem.
	Message string `json:"message" yaml:"message"`
	// Suggestion provides a suggested fix, if applicable.
	Suggestion string `json:"suggestion" yaml:"suggestion"`
	// Location provides the location of the problem. If both
	// `Location` and `Descriptor` are specified, the location
	// is then used from `Location` instead of `Descriptor`.
	Location Location `json:"location" yaml:"location"`
	// Descriptor provides the descriptor related
	// to the problem. If present and `Location` is not
	// specified, then the starting location of the descriptor
	// is used as the location of the problem.
	Descriptor protoreflect.Descriptor `json:"-" yaml:"-"`

	// The following fields will be set by Runtime.
	Category Category `json:"category" yaml:"category"`
	RuleID   RuleName `json:"rule_id" yaml:"rule_id"`
	FilePath string   `json:"file_path" yaml:"file_path"`
}
