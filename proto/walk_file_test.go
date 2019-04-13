package proto

import (
	"errors"
	"io/ioutil"
	"log"
	"path/filepath"
	"testing"

	"github.com/golang/protobuf/v2/proto"
	"github.com/golang/protobuf/v2/reflect/protodesc"
	"github.com/golang/protobuf/v2/reflect/protoreflect"
	descriptorpb "github.com/golang/protobuf/v2/types/descriptor"
	"github.com/jgeewax/api-linter/proto/mocks"
	"github.com/stretchr/testify/mock"
)

//go:generate protoc --include_source_info --descriptor_set_out=testdata/walk_file_test.protoset --proto_path=testdata testdata/walk_file_test.proto
//go:generate mockery -name Consumer

func TestWalkDescriptor(t *testing.T) {
	consumer := new(mocks.Consumer)
	f := readProtoFile("walk_file_test.protoset")

	// 15 = 1 file + 4 messages + 3 fields + 2 enums + 2 enum values + 1 oneof + 1 service + 1 method
	numDescriptors := 15
	consumer.On("Consume", mock.Anything).Return(nil).Times(numDescriptors)

	WalkFile(f, consumer)
	consumer.AssertExpectations(t)
}

func TestWalkDescriptorWithErr(t *testing.T) {
	consumer := new(mocks.Consumer)
	f := readProtoFile("walk_file_test.protoset")
	errStop := errors.New("stop")

	// just the file descriptor itself.
	numDescriptors := 1
	consumer.On("Consume", mock.Anything).Return(errStop).Times(numDescriptors)

	WalkFile(f, consumer)
	consumer.AssertExpectations(t)
}

func readProtoFile(fileName string) protoreflect.FileDescriptor {
	path := filepath.Join("testdata", fileName)
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Unable to open %s: %v", path, err)
	}
	protoset := &descriptorpb.FileDescriptorSet{}
	if err := proto.Unmarshal(bs, protoset); err != nil {
		log.Fatalf("Unable to parse %T from %s: %v", protoset, path, err)
	}
	f, err := protodesc.NewFile(protoset.GetFile()[0], nil)
	if err != nil {
		log.Fatalf("protodesc.NewFile() error: %v", err)
	}
	return f
}
