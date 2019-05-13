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
	"github.com/golang/protobuf/v2/reflect/protoreflect"
)

// Consumer represents an operation that consumes a single Descriptor.
type Consumer interface {
	Consume(protoreflect.Descriptor) error
}

// Walk travels in a Descriptor, such as FileDescriptor, MessageDescriptor, etc.
// The travel will continue to the nested types. For example, starting from a
// FileDescriptor, the visiting will continue to the nested Enum-, Extension-,
// Message-, and ServiceDescriptors. It will apply a Consumer to each encountered
// Descriptor until EOF or an error returned by the Consumer.
func Walk(d protoreflect.Descriptor, c Consumer) error {
	if err := c.Consume(d); err != nil {
		return err
	}

	// travel to the nested types.
	switch desc := d.(type) {
	case protoreflect.FileDescriptor:
		for i := 0; i < desc.Enums().Len(); i++ {
			if err := Walk(desc.Enums().Get(i), c); err != nil {
				return err
			}
		}
		for i := 0; i < desc.Extensions().Len(); i++ {
			if err := Walk(desc.Extensions().Get(i), c); err != nil {
				return err
			}
		}
		for i := 0; i < desc.Messages().Len(); i++ {
			if err := Walk(desc.Messages().Get(i), c); err != nil {
				return err
			}
		}
		for i := 0; i < desc.Services().Len(); i++ {
			if err := Walk(desc.Services().Get(i), c); err != nil {
				return err
			}
		}
	case protoreflect.MessageDescriptor:
		for i := 0; i < desc.Enums().Len(); i++ {
			if err := Walk(desc.Enums().Get(i), c); err != nil {
				return err
			}
		}
		for i := 0; i < desc.Extensions().Len(); i++ {
			if err := Walk(desc.Extensions().Get(i), c); err != nil {
				return err
			}
		}
		for i := 0; i < desc.Fields().Len(); i++ {
			if err := Walk(desc.Fields().Get(i), c); err != nil {
				return err
			}
		}
		for i := 0; i < desc.Messages().Len(); i++ {
			if err := Walk(desc.Messages().Get(i), c); err != nil {
				return err
			}
		}
		for i := 0; i < desc.Oneofs().Len(); i++ {
			if err := Walk(desc.Oneofs().Get(i), c); err != nil {
				return err
			}
		}
	case protoreflect.EnumDescriptor:
		for i := 0; i < desc.Values().Len(); i++ {
			if err := Walk(desc.Values().Get(i), c); err != nil {
				return err
			}
		}
	case protoreflect.ServiceDescriptor:
		for i := 0; i < desc.Methods().Len(); i++ {
			if err := Walk(desc.Methods().Get(i), c); err != nil {
				return err
			}
		}
	}

	return nil
}
