package lint

import "github.com/golang/protobuf/v2/reflect/protoreflect"

// Request defines input data for a rule to perform linting.
type Request struct {
	protoFile protoreflect.FileDescriptor

	descriptorSource DescriptorSource
}

func NewRequest(protoFile protoreflect.FileDescriptor, descriptorSource DescriptorSource) Request {
	return Request{
		protoFile: protoFile,
		descriptorSource: descriptorSource,
	}
}

// ProtoFile returns a `protoreflect.FileDescriptor`
//
// The `protoreflect.FileDescriptor` includes information about a protofile as well as methods to
// find elements that it contains.
func (r Request) ProtoFile() protoreflect.FileDescriptor {
	return r.protoFile
}

// DescriptorSource returns a `DescriptorSource`.
//
// The returned `DescriptorSource` contains additional information
// about a protobuf descriptor, such as comments and location lookups.
func (r Request) DescriptorSource() DescriptorSource {
	return r.descriptorSource
}
