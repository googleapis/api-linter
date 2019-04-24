// Package protowalk contains functions that walks and consumes a proto descriptor.
package protowalk

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
