package protovisit

import (
	"github.com/golang/protobuf/v2/reflect/protoreflect"
)

// DescriptorVisitor defines how to travel in a descriptor.
// See MessageVisitor for more details.
type DescriptorVisitor interface {
	// PreVisit will be invoked when a descriptor is encountered
	// before applying any visiting functions or traveling down
	// to its nested ones. If an error is returned, the entire
	// visiting stops -- expected that the special value ErrSkip
	// is returned, which indicates that this descriptor will be
	// skipped; or ErrSkipNested, which indicates that the nested
	// ones will be skipped; or ErrSkipVisiting, which indicates
	// that this descriptor will be skipped but visiting will
	// continue to the nested ones.
	PreVisit(protoreflect.Descriptor) error
	// PostVisit will be invoked in the end of visiting. If an
	// error is returned, the entire visiting stops.
	PostVisit(protoreflect.Descriptor) error
}

// DescriptorVisiting visits a descriptor.
type DescriptorVisiting interface {
	VisitDescriptor(protoreflect.Descriptor)
}

// WalkDescriptor uses the visitor to travel in the file,
// and applies the visiting functions on each encountered descriptor.
//
// Note: it does not visit the FileDescriptor itself.
func WalkDescriptor(f protoreflect.FileDescriptor, visitor DescriptorVisitor, visiting DescriptorVisiting) error {
	for i := 0; i < f.Enums().Len(); i++ {
		if err := walkDescriptor(f.Enums().Get(i), visitor, visiting); err != nil {
			return err
		}
	}
	for i := 0; i < f.Extensions().Len(); i++ {
		if err := walkDescriptor(f.Extensions().Get(i), visitor, visiting); err != nil {
			return err
		}
	}
	for i := 0; i < f.Messages().Len(); i++ {
		if err := walkDescriptor(f.Messages().Get(i), visitor, visiting); err != nil {
			return err
		}
	}
	for i := 0; i < f.Services().Len(); i++ {
		if err := walkDescriptor(f.Services().Get(i), visitor, visiting); err != nil {
			return err
		}
	}
	return nil
}

func walkDescriptor(d protoreflect.Descriptor, visitor DescriptorVisitor, visiting DescriptorVisiting) error {
	err := visitor.PreVisit(d)
	if err == ErrSkip {
		return nil
	}

	if err == nil || err == ErrSkipVisiting {
		switch d.(type) {
		case protoreflect.MessageDescriptor:
			m := d.(protoreflect.MessageDescriptor)
			for i := 0; i < m.Enums().Len(); i++ {
				if err := walkDescriptor(m.Enums().Get(i), visitor, visiting); err != nil {
					return err
				}
			}
			for i := 0; i < m.Extensions().Len(); i++ {
				if err := walkDescriptor(m.Enums().Get(i), visitor, visiting); err != nil {
					return err
				}
			}
			for i := 0; i < m.Fields().Len(); i++ {
				if err := walkDescriptor(m.Fields().Get(i), visitor, visiting); err != nil {
					return err
				}
			}
			for i := 0; i < m.Messages().Len(); i++ {
				if err := walkDescriptor(m.Messages().Get(i), visitor, visiting); err != nil {
					return err
				}
			}
			for i := 0; i < m.Oneofs().Len(); i++ {
				if err := walkDescriptor(m.Oneofs().Get(i), visitor, visiting); err != nil {
					return err
				}
			}
		case protoreflect.EnumDescriptor:
			e := d.(protoreflect.EnumDescriptor)
			for i := 0; i < e.Values().Len(); i++ {
				if err := walkDescriptor(e.Values().Get(i), visitor, visiting); err != nil {
					return err
				}
			}
		case protoreflect.ServiceDescriptor:
			s := d.(protoreflect.ServiceDescriptor)
			for i := 0; i < s.Methods().Len(); i++ {
				if err := walkDescriptor(s.Methods().Get(i), visitor, visiting); err != nil {
					return err
				}
			}
		}
	}

	if err == nil || err == ErrSkipNested {
		visiting.VisitDescriptor(d)
	}

	if err == nil || err == ErrSkipNested || err == ErrSkipVisiting {
		return visitor.PostVisit(d)
	}

	return err
}

// SimpleDescriptorVisitor visits all descriptors in a file.
type SimpleDescriptorVisitor struct{}

// PreVisit does nothing and return nil.
func (v SimpleDescriptorVisitor) PreVisit(protoreflect.Descriptor) error { return nil }

// PostVisit does nothing and return nil.
func (v SimpleDescriptorVisitor) PostVisit(protoreflect.Descriptor) error { return nil }
