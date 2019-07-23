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

package lint

import (
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
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
//
// Deprecated: Use NewRequestFromFileDescriptor instead.
func NewProtoRequest(fd *descriptorpb.FileDescriptorProto, reg *protoregistry.Files) (Request, error) {
	f, err := protodesc.NewFile(fd, reg)
	if err != nil {
		return Request{}, err
	}
	return NewRequestFromFileDescriptor(f)
}

// NewRequestFromFileDescriptor creates a linting Request from the given FileDescriptor.
func NewRequestFromFileDescriptor(f protoreflect.FileDescriptor) (Request, error) {
	s, err := newDescriptorSource(f)
	return Request{
		fileDesc:   f,
		descSource: s,
	}, err
}
