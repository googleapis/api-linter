// Package protovisit defines visitors for top-level descriptors in a protobuf file.
// E.g., message, enum, service, etc.
//
// For a top-level descriptor, it contains two interfaces:
//   1. Visitor, e.g., MessageVisitor, that defines how to travel in the descriptor;
//   2. Visiting, e.g., MessageVisiting, that defines what functions to be applied on the descriptor and its elements;
// and a work function, that uses the visitor to travel applies the visiting function on the encountered descriptors;
package protovisit

import (
	"errors"

	"github.com/golang/protobuf/v2/reflect/protoreflect"
)

// ErrSkip indicates that the encountered descriptor is to be skipped.
var ErrSkip = errors.New("skip this descriptor")

// ErrSkipVisiting indicates not to apply visiting functions on the descriptor,
// but the traverse will continue to the subtree and PostVisit will be invoked in the end.
var ErrSkipVisiting = errors.New("skip visiting this descriptor, but keep going to its subtree")

// ErrSkipNested indicates that the nested messages will be skipped.
var ErrSkipNested = errors.New("skip the nested messages")

// EnumVisitor defines how to travel in an enum.
// See MessageVisitor for more detail.
type EnumVisitor interface {
	PreVisit(protoreflect.EnumDescriptor) error
	PostVisit(protoreflect.EnumDescriptor) error
}

// EnumVisiting defines a collection of functions that can be applied to an enum and
// its values.
type EnumVisiting interface {
	VisitEnum(protoreflect.EnumDescriptor)
	VisitEnumValue(protoreflect.EnumValueDescriptor)
}

// WalkEnum uses the visitor to travel in the top-level enums and applies the visiting functions to
// each encountered enum.
func WalkEnum(f protoreflect.FileDescriptor, visitor EnumVisitor, visiting EnumVisiting) error {
	for i := 0; i < f.Enums().Len(); i++ {
		if err := walkEnum(f.Enums().Get(i), visitor, visiting); err != nil {
			return err
		}
	}
	return nil
}

func walkEnum(e protoreflect.EnumDescriptor, visitor EnumVisitor, visiting EnumVisiting) error {
	err := visitor.PreVisit(e)
	if err == ErrSkip {
		return nil
	}

	if err == nil {
		visitEnum(e, visiting)
	}

	if err == nil || err == ErrSkipVisiting {
		return visitor.PostVisit(e)
	}
	return err
}

func visitEnum(e protoreflect.EnumDescriptor, v EnumVisiting) {
	v.VisitEnum(e)

	for i := 0; i < e.Values().Len(); i++ {
		v.VisitEnumValue(e.Values().Get(i))
	}
}

// ExtensionVisitor defines how to travel in an extension (field).
// See MessageVisitor for more details.
type ExtensionVisitor interface {
	PreVisit(protoreflect.ExtensionDescriptor) error
	PostVisit(protoreflect.ExtensionDescriptor) error
}

// ExtensionVisiting defines a collection of functions that can be applied to
// a extension (field).
type ExtensionVisiting interface {
	VisitExtension(protoreflect.ExtensionDescriptor)
}

// WalkExtension uses the visitor to travel in top-level extensions and applies the visiting functions
// to each encountered extension.
func WalkExtension(f protoreflect.FileDescriptor, visitor ExtensionVisitor, visiting ExtensionVisiting) error {
	for i := 0; i < f.Extensions().Len(); i++ {
		if err := walkExtension(f.Extensions().Get(i), visitor, visiting); err != nil {
			return err
		}
	}
	return nil
}

func walkExtension(e protoreflect.ExtensionDescriptor, visitor ExtensionVisitor, visiting ExtensionVisiting) error {
	err := visitor.PreVisit(e)
	if err == ErrSkip {
		return nil
	}

	if err == nil {
		visitExtension(e, visiting)
	}

	if err == nil || err == ErrSkipVisiting {
		return visitor.PostVisit(e)
	}
	return err
}

func visitExtension(e protoreflect.ExtensionDescriptor, visiting ExtensionVisiting) {
	visiting.VisitExtension(e)
}

// MessageVisitor defines how to travel in a message.
type MessageVisitor interface {
	// PreVisit will be invoked when the visitor enters the message,
	// before applying visiting functions or traveling down to its
	// nested ones. If an error is returned, the whole visiting stops --
	// except that the special value ErrSkip is returned, which
	// indicates that this message is to be skipped, or ErrSkipNested,
	// which indicates that the nested ones will be skipped, or ErrSkipVisiting
	// which indicates that this message will be skipped visiting but visiting
	// will continue to its nested messages -- in all cases, the processing will
	// continue to the rest, same-level (sibling) messages.
	PreVisit(protoreflect.MessageDescriptor) error

	// PostVisit will be invoked in the end of visiting. If an error is returned,
	// the whole visiting stops.
	PostVisit(protoreflect.MessageDescriptor) error
}

// MessageVisiting defines a collection of functions that can be applied to a message and
// its elements.
type MessageVisiting interface {
	EnumVisiting
	ExtensionVisiting
	VisitMessage(protoreflect.MessageDescriptor)
	VisitField(protoreflect.FieldDescriptor)
	VisitOneof(protoreflect.OneofDescriptor)
}

// WalkMessage uses the visitor to travel in top-level messages and applies the visiting functions on each
// encountered message.
func WalkMessage(file protoreflect.FileDescriptor, visitor MessageVisitor, visiting MessageVisiting) error {
	for i := 0; i < file.Messages().Len(); i++ {
		if err := walkMessage(file.Messages().Get(i), visitor, visiting); err != nil {
			return err
		}
	}
	return nil
}

// depth-first walking in a message.
func walkMessage(msg protoreflect.MessageDescriptor, visitor MessageVisitor, visiting MessageVisiting) error {
	err := visitor.PreVisit(msg)
	if err == ErrSkip {
		return nil
	}

	if err == nil || err == ErrSkipVisiting {
		nested := msg.Messages()
		for i := 0; i < nested.Len(); i++ {
			if err := walkMessage(nested.Get(i), visitor, visiting); err != nil {
				return err
			}
		}
	}

	if err == nil || err == ErrSkipNested {
		visitMessage(msg, visiting)
	}

	if err == nil || err == ErrSkipNested || err == ErrSkipVisiting {
		return visitor.PostVisit(msg)
	}

	return err
}

func visitMessage(msg protoreflect.MessageDescriptor, v MessageVisiting) {
	v.VisitMessage(msg)

	enums := msg.Enums()
	for i := 0; i < enums.Len(); i++ {
		visitEnum(enums.Get(i), v)
	}

	extensions := msg.Extensions()
	for i := 0; i < extensions.Len(); i++ {
		visitExtension(extensions.Get(i), v)
	}

	fields := msg.Fields()
	for i := 0; i < fields.Len(); i++ {
		v.VisitField(fields.Get(i))
	}

	oneofs := msg.Oneofs()
	for i := 0; i < oneofs.Len(); i++ {
		v.VisitOneof(oneofs.Get(i))
	}
}

// ServiceVisitor defines how to travel in a service.
// See MessageVisitor for more detail.
type ServiceVisitor interface {
	PreVisit(protoreflect.ServiceDescriptor) error

	PostVisit(protoreflect.ServiceDescriptor) error
}

// ServiceVisiting defines a collection of functions that can be applied to a service and
// its elements.
type ServiceVisiting interface {
	VisitService(protoreflect.ServiceDescriptor)
	VisitMethod(protoreflect.MethodDescriptor)
}

// WalkService uses the visitor to travel in top-level services and applies the visiting functions on the
// encountered service.
func WalkService(f protoreflect.FileDescriptor, visitor ServiceVisitor, visiting ServiceVisiting) error {
	for i := 0; i < f.Services().Len(); i++ {
		if err := walkService(f.Services().Get(i), visitor, visiting); err != nil {
			return err
		}
	}
	return nil
}

func walkService(s protoreflect.ServiceDescriptor, visitor ServiceVisitor, visiting ServiceVisiting) error {
	err := visitor.PreVisit(s)
	if err == ErrSkip {
		return nil
	}

	if err == nil {
		visitService(s, visiting)
	}

	if err == nil || err == ErrSkipVisiting {
		return visitor.PostVisit(s)
	}
	return err
}

// visitService applies the visiting functions on the service and its elements.
func visitService(s protoreflect.ServiceDescriptor, v ServiceVisiting) {
	v.VisitService(s)

	for i := 0; i < s.Methods().Len(); i++ {
		v.VisitMethod(s.Methods().Get(i))
	}

}

// SimpleEnumVisitor visits every top-level enums in the file.
type SimpleEnumVisitor struct{}

// PreVisit does nothing.
func (v SimpleEnumVisitor) PreVisit(m protoreflect.EnumDescriptor) error { return nil }

// PostVisit does nothing.
func (v SimpleEnumVisitor) PostVisit(m protoreflect.EnumDescriptor) error { return nil }

// SimpleExtensionVisitor visits every top-level extensions in the file.
type SimpleExtensionVisitor struct{}

// PreVisit does nothing.
func (v SimpleExtensionVisitor) PreVisit(m protoreflect.ExtensionDescriptor) error { return nil }

// PostVisit does nothing.
func (v SimpleExtensionVisitor) PostVisit(m protoreflect.ExtensionDescriptor) error { return nil }

// SimpleMessageVisitor visits every message it encounters in the file.
type SimpleMessageVisitor struct{}

// PreVisit does nothing.
func (v SimpleMessageVisitor) PreVisit(m protoreflect.MessageDescriptor) error { return nil }

// PostVisit does nothing.
func (v SimpleMessageVisitor) PostVisit(m protoreflect.MessageDescriptor) error { return nil }

// SimpleServiceVisitor implements ServiceVisitor that visits all services.
type SimpleServiceVisitor struct{}

// PreVisit does nothing.
func (s SimpleServiceVisitor) PreVisit(protoreflect.ServiceDescriptor) error { return nil }

// PostVisit does nothing.
func (s SimpleServiceVisitor) PostVisit(protoreflect.ServiceDescriptor) error { return nil }

// ScopedMessageVisitor calls `Predicate` to determine if an encountered
// message is in scope for visiting. If yes, the message and its subtree
// will be subject to visiting.
type ScopedMessageVisitor struct {
	Predicate func(protoreflect.MessageDescriptor) bool

	depth int
}

func (v *ScopedMessageVisitor) isInScope() bool {
	return v.depth > 0
}

// PreVisit calls `Predicate` to determine if the message is in scope for visiting.
// If the `Predicate` is nil, the message will be always in scope.
func (v *ScopedMessageVisitor) PreVisit(m protoreflect.MessageDescriptor) error {
	if v.Predicate == nil || v.isInScope() || v.Predicate(m) {
		v.depth++
	}
	if v.isInScope() {
		return nil
	}
	return ErrSkipVisiting
}

// PostVisit reduces the depth of scope when leaving the message.
func (v *ScopedMessageVisitor) PostVisit(m protoreflect.MessageDescriptor) error {
	if v.isInScope() {
		v.depth--
	}
	return nil
}
