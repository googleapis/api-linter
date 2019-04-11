package lint

import "github.com/golang/protobuf/v2/reflect/protoreflect"

// Request defines input data for a rule to perform linting.
type Request struct {
	protoFile  protoreflect.FileDescriptor
	descSource DescriptorSource
}

// NewProtoFileRequest creates a linting request for a .proto file.
func NewProtoFileRequest(protoFile protoreflect.FileDescriptor, descSource DescriptorSource) Request {
	return Request{
		protoFile:  protoFile,
		descSource: descSource,
	}
}
