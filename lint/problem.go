package lint

import (
	"github.com/golang/protobuf/v2/reflect/protoreflect"
)

// Problem contains information about a result produced by an API Linter.
type Problem struct {
	// Message provides a short description of the problem.
	Message string
	// Suggestion provides a suggested fix, if applicable.
	Suggestion string
	// Location provides the location of the problem. If both
	// `Location` and `Descriptor` are specified, the location
	// is then used from `Location` instead of `Descriptor`.
	Location Location
	// Descriptor provides the protobuf descriptor related
	// to the problem. If present and `Location` is not
	// specified, then the starting location of the descriptor
	// is used as the location of the problem.
	Descriptor protoreflect.Descriptor
}
