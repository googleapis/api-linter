package rules

import (
	"github.com/golang/protobuf/v2/proto"
	descriptorpb "github.com/golang/protobuf/v2/types/descriptor"
	"log"
)

func protoDescriptorProtoFromSource(source []byte) *descriptorpb.FileDescriptorProto {
	protoset := &descriptorpb.FileDescriptorSet{}
	if err := proto.Unmarshal(source, protoset); err != nil {
		log.Fatalf("Unable to parse %T from source: %v.", protoset, err)
	}

	return protoset.GetFile()[0]
}
