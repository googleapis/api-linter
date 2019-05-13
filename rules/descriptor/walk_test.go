// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package descriptor

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
)

//go:generate protoc --include_source_info --descriptor_set_out=testdata/test.protoset --proto_path=testdata testdata/test.proto

type mockConsumer struct {
	count int
	err   error
}

func (m *mockConsumer) Consume(d protoreflect.Descriptor) error {
	m.count++
	return m.err
}

func TestWalk(t *testing.T) {
	f := readProtoFile("test.protoset")
	tests := []struct {
		descriptor protoreflect.Descriptor
		num        int
	}{
		{
			descriptor: f,
			num:        18, // 18 = 1 file + 5 messages + 3 fields + 2 enums + 2 enum values + 1 oneof + 1 service + 1 method + 2 extensions.
		},
		{
			descriptor: f.Enums().Get(0),
			num:        2, // 2 = 1 enum + 1 value
		},
		{
			descriptor: f.Messages().Get(0),
			num:        8, // 8 = 2 messages + 1 enum + 1 value + 1 oneof + 3 fields
		},
		{
			descriptor: f.Services().Get(0),
			num:        2, // 2 = 1 service + 1 method
		},
		{
			descriptor: f.Messages().Get(0).Fields().Get(0),
			num:        1, // 1 = 1 field
		},
		{
			descriptor: f.Messages().Get(0).Oneofs().Get(0),
			num:        1, // 1 = 1 oneof
		},
	}
	for _, test := range tests {
		consumer := new(mockConsumer)

		Walk(test.descriptor, consumer)
		if consumer.count != test.num {
			t.Errorf("Walk(%s): Got %d desriptors, but wanted %d", test.descriptor.FullName(), consumer.count, test.num)
		}
	}
}

func TestWalkWithErr(t *testing.T) {
	consumer := &mockConsumer{
		err: errors.New("stop"),
	}
	f := readProtoFile("test.protoset")

	Walk(f, consumer)
	if consumer.count != 1 {
		t.Errorf("Walk(%s) with error: got %d descriptors, but wanted 1", f.FullName(), consumer.count)
	}
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
