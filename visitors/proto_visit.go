// Package visitors contains Protobuf file visitors.
package visitors

import (
	"errors"

	"github.com/golang/protobuf/reflect/protoreflect"
)

// ErrSkip indicates that the encountered message is to be skipped.
var ErrSkip = errors.New("skip this message")

// ErrSkipSubtree indicates that the sub-messages of the encountered message
// are to be skipped.
var ErrSkipSubtree = errors.New("skip the sub-messages")

// MessageVisitor defines a visitor that travels in
// a protobuf message.
type MessageVisitor interface {
	// PreVisit will be invoked when the visitor enters the message,
	// before applying visiting functions or traveling down to its
	// nested ones. If an error is returned, the whole visiting stops --
	// except that the special value ErrSkip is returned, which
	// indicates that this message is to be skipped, or ErrSkipSubtree,
	// which indicates that the nested ones will be skipped; in
	// both cases, the processing will continue to the rest,
	// same-level messages.
	PreVisit(protoreflect.MessageDescriptor) error

	// Visit will be invoked only when `PreVisit` returns no error
	// or `ErrSkipSubtree`.
	Visit(protoreflect.MessageDescriptor)

	// PostVisit will be invoked after `Visit`. If an error is returned,
	// the whole visiting stops.
	PostVisit(protoreflect.MessageDescriptor) error
}

// WalkMessage travels in a `FileDescriptor`, and applies
// the `MessageVisitor` when a message is encountered.
func WalkMessage(file protoreflect.FileDescriptor, visitor MessageVisitor) error {
	for i := 0; i < file.Messages().Len(); i++ {
		if err := walkMessage(file.Messages().Get(i), visitor); err != nil {
			return err
		}
	}
	return nil
}

// depth-first walking in a message.
func walkMessage(msg protoreflect.MessageDescriptor, visitor MessageVisitor) error {
	err := visitor.PreVisit(msg)
	if err == ErrSkip {
		return nil
	}

	if err == nil {
		nested := msg.Messages()
		for i := 0; i < nested.Len(); i++ {
			if err := walkMessage(nested.Get(i), visitor); err != nil {
				return err
			}
		}
	}

	if err == nil || err == ErrSkipSubtree {
		visitor.Visit(msg)
		return visitor.PostVisit(msg)
	}

	return err
}

// MessageVisitingFuncs contains functions for visiting name,
// fields, enums, and oneofs of a message.
type MessageVisitingFuncs struct {
	MessageVisit   func(protoreflect.MessageDescriptor)
	FieldVisit     func(protoreflect.FieldDescriptor)
	EnumVisit      func(protoreflect.EnumDescriptor)
	EnumValueVisit func(protoreflect.EnumValueDescriptor)
	OneofVisit     func(protoreflect.OneofDescriptor)
}

// visitMessage applies `MessageVisit`, if exists, on the message.
func (v MessageVisitingFuncs) visitMessage(m protoreflect.MessageDescriptor) {
	if v.MessageVisit != nil {
		v.MessageVisit(m)
	}
}

// visitField applies `FieldVisit`, if exists, on a message field.
func (v MessageVisitingFuncs) visitField(f protoreflect.FieldDescriptor) {
	if v.FieldVisit != nil {
		v.FieldVisit(f)
	}
}

// visitEnum applies `EnumVisit`, if exists, on a message enum.
func (v MessageVisitingFuncs) visitEnum(e protoreflect.EnumDescriptor) {
	if v.EnumVisit != nil {
		v.EnumVisit(e)
	}
	if v.EnumValueVisit != nil {
		for i := 0; i < e.Values().Len(); i++ {
			v.EnumValueVisit(e.Values().Get(i))
		}
	}
}

// visitOneof applies `OneofVisit`, if exists, on a message oneof.
func (v MessageVisitingFuncs) visitOneof(o protoreflect.OneofDescriptor) {
	if v.OneofVisit != nil {
		v.OneofVisit(o)
	}
}

func visitMessage(msg protoreflect.MessageDescriptor, v MessageVisitingFuncs) {
	v.visitMessage(msg)

	fields := msg.Fields()
	for i := 0; i < fields.Len(); i++ {
		v.visitField(fields.Get(i))
	}

	oneofs := msg.Oneofs()
	for i := 0; i < oneofs.Len(); i++ {
		v.visitOneof(oneofs.Get(i))
	}

	enums := msg.Enums()
	for i := 0; i < enums.Len(); i++ {
		v.visitEnum(enums.Get(i))
	}
}

// ScopedMessageVisitor calls `Predicate` to determine if an encountered
// message is in scope for visiting. If yes, the message and its subtree
// will be subject to visiting.
type ScopedMessageVisitor struct {
	Predicate func(protoreflect.MessageDescriptor) bool
	Funcs     MessageVisitingFuncs

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
	return nil
}

// Visit visits a message, and applies visiting `Funcs` on it, if the message is
// in scope.
func (v *ScopedMessageVisitor) Visit(m protoreflect.MessageDescriptor) {
	if v.isInScope() {
		visitMessage(m, v.Funcs)
	}
}

// PostVisit reduces the depth of scope when leaving the message.
func (v *ScopedMessageVisitor) PostVisit(m protoreflect.MessageDescriptor) error {
	if v.isInScope() {
		v.depth--
	}
	return nil
}

// SimpleMessageVisitor visits every message it encounters in the file,
// and applies the visiting functions.
type SimpleMessageVisitor struct {
	Funcs MessageVisitingFuncs
}

// PreVisit does nothing.
func (v *SimpleMessageVisitor) PreVisit(m protoreflect.MessageDescriptor) error {
	return nil
}

// Visit visits a message, and applies `Funcs` on it.
func (v *SimpleMessageVisitor) Visit(m protoreflect.MessageDescriptor) {
	visitMessage(m, v.Funcs)
}

// PostVisit does nothing.
func (v *SimpleMessageVisitor) PostVisit(m protoreflect.MessageDescriptor) error {
	return nil
}

// ServiceVisitor defines a visitor that can travel in a service.
type ServiceVisitor interface {
	// PreVisit will be invoked when the visitor enters the service,
	// before applying visiting functions the methods.
	// If an error is returned, the whole visiting stops --
	// except that the special value ErrSkip is returned, which
	// indicates that this service is to be skipped.
	PreVisit(protoreflect.ServiceDescriptor) error

	// Visit will be invoked only when `PreVisit` returns no error
	// or `ErrSkipSubtree`.
	Visit(protoreflect.ServiceDescriptor)

	// PostVisit will be invoked after `Visit`. If an error is returned,
	// the whole visiting stops.
	PostVisit(protoreflect.ServiceDescriptor) error
}

// WalkService travels in a `FileDescriptor`, and applies
// the `ServiceVisitor` when a service is encountered.
func WalkService(f protoreflect.FileDescriptor, v ServiceVisitor) error {
	for i := 0; i < f.Services().Len(); i++ {
		if err := walkService(f.Services().Get(i), v); err != nil {
			return err
		}
	}
	return nil
}

func walkService(s protoreflect.ServiceDescriptor, v ServiceVisitor) error {
	err := v.PreVisit(s)
	if err == ErrSkip {
		return nil
	}

	if err != nil {
		return err
	}

	if err == nil {
		v.Visit(s)
	}
	return v.PostVisit(s)
}

// ServiceVisitFuncs defines a set of functions that can visit a service.
type ServiceVisitFuncs struct {
	ServiceVisit func(protoreflect.ServiceDescriptor)
	MethodVisit  func(protoreflect.MethodDescriptor)
}

func visitService(s protoreflect.ServiceDescriptor, v ServiceVisitFuncs) {
	if v.ServiceVisit != nil {
		v.ServiceVisit(s)
	}
	if v.MethodVisit != nil {
		for i := 0; i < s.Methods().Len(); i++ {
			v.MethodVisit(s.Methods().Get(i))
		}
	}
}

// SimpleServiceVisitor implements ServiceVisitor that visits all services.
type SimpleServiceVisitor struct {
	Funcs ServiceVisitFuncs
}

// PreVisit does nothing.
func (s SimpleServiceVisitor) PreVisit(protoreflect.ServiceDescriptor) error {
	return nil
}

// Visit visits a service.
func (s SimpleServiceVisitor) Visit(d protoreflect.ServiceDescriptor) {
	visitService(d, s.Funcs)
}

// PostVisit does nothing.
func (s SimpleServiceVisitor) PostVisit(protoreflect.ServiceDescriptor) error {
	return nil
}
