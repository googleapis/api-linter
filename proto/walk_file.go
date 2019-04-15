// Package proto contains helper functions for visiting a .proto file.
package proto

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
	return walkDescriptor(d, c)
}

func walkDescriptor(d protoreflect.Descriptor, c Consumer) error {
	if err := c.Consume(d); err != nil {
		return err
	}

	// travel to the nested types.
	switch d.(type) {
	case protoreflect.FileDescriptor:
		f := d.(protoreflect.FileDescriptor)
		for i := 0; i < f.Enums().Len(); i++ {
			if err := walkDescriptor(f.Enums().Get(i), c); err != nil {
				return err
			}
		}
		for i := 0; i < f.Extensions().Len(); i++ {
			if err := walkDescriptor(f.Extensions().Get(i), c); err != nil {
				return err
			}
		}
		for i := 0; i < f.Messages().Len(); i++ {
			if err := walkDescriptor(f.Messages().Get(i), c); err != nil {
				return err
			}
		}
		for i := 0; i < f.Services().Len(); i++ {
			if err := walkDescriptor(f.Services().Get(i), c); err != nil {
				return err
			}
		}
	case protoreflect.MessageDescriptor:
		m := d.(protoreflect.MessageDescriptor)
		for i := 0; i < m.Enums().Len(); i++ {
			if err := walkDescriptor(m.Enums().Get(i), c); err != nil {
				return err
			}
		}
		for i := 0; i < m.Extensions().Len(); i++ {
			if err := walkDescriptor(m.Extensions().Get(i), c); err != nil {
				return err
			}
		}
		for i := 0; i < m.Fields().Len(); i++ {
			if err := walkDescriptor(m.Fields().Get(i), c); err != nil {
				return err
			}
		}
		for i := 0; i < m.Messages().Len(); i++ {
			if err := walkDescriptor(m.Messages().Get(i), c); err != nil {
				return err
			}
		}
		for i := 0; i < m.Oneofs().Len(); i++ {
			if err := walkDescriptor(m.Oneofs().Get(i), c); err != nil {
				return err
			}
		}
	case protoreflect.EnumDescriptor:
		e := d.(protoreflect.EnumDescriptor)
		for i := 0; i < e.Values().Len(); i++ {
			if err := walkDescriptor(e.Values().Get(i), c); err != nil {
				return err
			}
		}
	case protoreflect.ServiceDescriptor:
		s := d.(protoreflect.ServiceDescriptor)
		for i := 0; i < s.Methods().Len(); i++ {
			if err := walkDescriptor(s.Methods().Get(i), c); err != nil {
				return err
			}
		}
	}

	return nil
}
