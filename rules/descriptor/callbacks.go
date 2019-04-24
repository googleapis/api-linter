package descriptor

import (
	"github.com/golang/protobuf/v2/reflect/protoreflect"
	"github.com/jgeewax/api-linter/lint"
)

// Callbacks contains a collection of functions that will be called
// for every encountered, corresponding descriptor.
//
// Note: `DescriptorCallback`, if exists, is called only when no *specific*
// callback is available, i.e., it is a fallback.
type Callbacks struct {
	EnumCallback       func(protoreflect.EnumDescriptor, lint.DescriptorSource) ([]lint.Problem, error)
	EnumValueCallback  func(protoreflect.EnumValueDescriptor, lint.DescriptorSource) ([]lint.Problem, error)
	FieldCallback      func(protoreflect.FieldDescriptor, lint.DescriptorSource) ([]lint.Problem, error)
	ExtensionCallback  func(protoreflect.ExtensionDescriptor, lint.DescriptorSource) ([]lint.Problem, error)
	FileCallback       func(protoreflect.FileDescriptor, lint.DescriptorSource) ([]lint.Problem, error)
	MessageCallback    func(protoreflect.MessageDescriptor, lint.DescriptorSource) ([]lint.Problem, error)
	MethodCallback     func(protoreflect.MethodDescriptor, lint.DescriptorSource) ([]lint.Problem, error)
	ServiceCallback    func(protoreflect.ServiceDescriptor, lint.DescriptorSource) ([]lint.Problem, error)
	OneofCallback      func(protoreflect.OneofDescriptor, lint.DescriptorSource) ([]lint.Problem, error)
	DescriptorCallback func(protoreflect.Descriptor, lint.DescriptorSource) ([]lint.Problem, error)
}

// Apply invokes a specific, corresponding callback for the descriptor with the source.
// Only when no specific callback is available, the `DescriptorCallback` will be tried.
func (c Callbacks) Apply(d protoreflect.Descriptor, src lint.DescriptorSource) ([]lint.Problem, error) {
	switch desc := d.(type) {
	case protoreflect.EnumDescriptor:
		if c.EnumCallback != nil {
			return c.EnumCallback(desc, src)
		}
	case protoreflect.EnumValueDescriptor:
		if c.EnumValueCallback != nil {
			return c.EnumValueCallback(desc, src)
		}
	case protoreflect.FieldDescriptor:
		if desc.ExtendedType() != nil && c.ExtensionCallback != nil {
			return c.ExtensionCallback(desc, src)
		}
		if c.FieldCallback != nil {
			return c.FieldCallback(desc, src)
		}
	case protoreflect.FileDescriptor:
		if c.FileCallback != nil {
			return c.FileCallback(desc, src)
		}
	case protoreflect.MethodDescriptor:
		if c.MethodCallback != nil {
			return c.MethodCallback(desc, src)
		}
	case protoreflect.MessageDescriptor:
		if c.MessageCallback != nil {
			return c.MessageCallback(desc, src)
		}
	case protoreflect.ServiceDescriptor:
		if c.ServiceCallback != nil {
			return c.ServiceCallback(desc, src)
		}
	case protoreflect.OneofDescriptor:
		if c.OneofCallback != nil {
			return c.OneofCallback(desc, src)
		}
	}

	// fallback to the general callback.
	if c.DescriptorCallback != nil {
		return c.DescriptorCallback(d, src)
	}

	return nil, nil
}
