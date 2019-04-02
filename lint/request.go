package lint

import "github.com/golang/protobuf/v2/reflect/protoreflect"

// Request defines input data for a rule to perform linting.
type Request interface {
	// ProtoFile returns a `FileDescriptor` when the rule's `FileTypes`
	// contains `ProtoFile`.
	ProtoFile() protoreflect.FileDescriptor
	Context() Context
}
