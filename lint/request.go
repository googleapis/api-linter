package lint

import (
	"github.com/golang/protobuf/v2/reflect/protodesc"
	"github.com/golang/protobuf/v2/reflect/protoreflect"
	descriptorpb "github.com/golang/protobuf/v2/types/descriptor"
)

// Request defines input data for a rule to perform linting.
type Request struct {
	fileDesc   protoreflect.FileDescriptor
	descSource DescriptorSource
}

// ProtoFile returns a FileDescriptor of the .proto file that will be linted.
func (r Request) ProtoFile() protoreflect.FileDescriptor {
	return r.fileDesc
}

// DescriptorSource returns a DescriptorSource that contains additional source
// information for the .proto file that will be linted.
func (r Request) DescriptorSource() DescriptorSource {
	return r.descSource
}

// NewProtoRequest creates a linting Request for a .proto file.
func NewProtoRequest(fd *descriptorpb.FileDescriptorProto) (Request, error) {
	f, err := protodesc.NewFile(fd, nil)
	if err != nil {
		return Request{}, err
	}
	s, err := newDescriptorSource(fd)
	return Request{
		fileDesc:   f,
		descSource: s,
	}, err
}
