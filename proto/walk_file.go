// Package proto contains helper functions for visiting a .proto fildesc.
package proto

import (
	"github.com/golang/protobuf/v2/reflect/protoreflect"
)

// DescriptorConsumer represents an operation that consumes a single Descriptor.
type DescriptorConsumer interface {
	ConsumeDescriptor(protoreflect.Descriptor) error
}

// WalkDescriptor travels in a Descriptor, such as FileDescriptor, MessageDescriptor, etc.
// The travel will continue to the nested types. For example, starting from a
// FileDescriptor, the visiting will continue to the nested Enum-, Extension-,
// Message-, and ServiceDescriptors. It will apply a DescriptorConsumer to each encountered
// Descriptor until EOF or an error returned by the DescriptorConsumer.
func WalkDescriptor(d protoreflect.Descriptor, c DescriptorConsumer) error {
	return walkDescriptor(d, c)
}

func walkDescriptor(d protoreflect.Descriptor, c DescriptorConsumer) error {
	if err := c.ConsumeDescriptor(d); err != nil {
		return err
	}

	// travel to the nested types.
	switch desc := d.(type) {
	case protoreflect.FileDescriptor:
		for i := 0; i < desc.Enums().Len(); i++ {
			if err := walkDescriptor(desc.Enums().Get(i), c); err != nil {
				return err
			}
		}
		for i := 0; i < desc.Extensions().Len(); i++ {
			if err := walkDescriptor(desc.Extensions().Get(i), c); err != nil {
				return err
			}
		}
		for i := 0; i < desc.Messages().Len(); i++ {
			if err := walkDescriptor(desc.Messages().Get(i), c); err != nil {
				return err
			}
		}
		for i := 0; i < desc.Services().Len(); i++ {
			if err := walkDescriptor(desc.Services().Get(i), c); err != nil {
				return err
			}
		}
	case protoreflect.MessageDescriptor:
		for i := 0; i < desc.Enums().Len(); i++ {
			if err := walkDescriptor(desc.Enums().Get(i), c); err != nil {
				return err
			}
		}
		for i := 0; i < desc.Extensions().Len(); i++ {
			if err := walkDescriptor(desc.Extensions().Get(i), c); err != nil {
				return err
			}
		}
		for i := 0; i < desc.Fields().Len(); i++ {
			if err := walkDescriptor(desc.Fields().Get(i), c); err != nil {
				return err
			}
		}
		for i := 0; i < desc.Messages().Len(); i++ {
			if err := walkDescriptor(desc.Messages().Get(i), c); err != nil {
				return err
			}
		}
		for i := 0; i < desc.Oneofs().Len(); i++ {
			if err := walkDescriptor(desc.Oneofs().Get(i), c); err != nil {
				return err
			}
		}
	case protoreflect.EnumDescriptor:
		for i := 0; i < desc.Values().Len(); i++ {
			if err := walkDescriptor(desc.Values().Get(i), c); err != nil {
				return err
			}
		}
	case protoreflect.ServiceDescriptor:
		for i := 0; i < desc.Methods().Len(); i++ {
			if err := walkDescriptor(desc.Methods().Get(i), c); err != nil {
				return err
			}
		}
	}

	return nil
}
