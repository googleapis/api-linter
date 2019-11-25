package descrule

import (
	"github.com/googleapis/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
)

// FileRule defines a lint rule that checks a file as a whole.
type FileRule struct {
	RuleName lint.RuleName

	// LintFile accepts a FileDescriptor and lints it, returning a slice of
	// lint.Problems it finds.
	LintFile func(*desc.FileDescriptor) []lint.Problem

	// OnlyIf accepts a FileDescriptor and determines whether this rule
	// is applicable.
	OnlyIf func(*desc.FileDescriptor) bool

	noPositional struct{}
}

// Name returns the name of the rule.
func (r *FileRule) Name() lint.RuleName {
	return r.RuleName
}

// Lint accepts a Descriptor and applies LintFile on it
// only if the descriptor is a FileDescriptor.
func (r *FileRule) Lint(d lint.Descriptor) []lint.Problem {
	if f, ok := d.(*File); ok {
		if r.OnlyIf == nil || r.OnlyIf(f.descriptor) {
			return r.LintFile(f.descriptor)
		}
	}
	return nil
}

// MessageRule defines a lint rule that is run on each message in the file.
//
// Both top-level messages and nested messages are visited.
type MessageRule struct {
	RuleName lint.RuleName

	// LintMessage accepts a MessageDescriptor and lints it, returning a slice
	// of lint.Problems it finds.
	LintMessage func(*desc.MessageDescriptor) []lint.Problem

	// OnlyIf accepts a MessageDescriptor and determines whether this rule
	// is applicable.
	OnlyIf func(*desc.MessageDescriptor) bool

	noPositional struct{}
}

// Name returns the name of the rule.
func (r *MessageRule) Name() lint.RuleName {
	return r.RuleName
}

// Lint accepts a Descriptor and applies LintMessage on it
// only if the descriptor is a MessageDescriptor.
func (r *MessageRule) Lint(d lint.Descriptor) []lint.Problem {
	if m, ok := d.(*Message); ok {
		if r.OnlyIf == nil || r.OnlyIf(m.descriptor) {
			return r.LintMessage(m.descriptor)
		}
	}
	return nil
}

// FieldRule defines a lint rule that is run on each field within a file.
type FieldRule struct {
	RuleName lint.RuleName

	// LintField accepts a FieldDescriptor and lints it, returning a slice of
	// lint.Problems it finds.
	LintField func(*desc.FieldDescriptor) []lint.Problem

	// OnlyIf accepts a FieldDescriptor and determines whether this rule
	// is applicable.
	OnlyIf func(*desc.FieldDescriptor) bool

	noPositional struct{}
}

// Name returns the name of the rule.
func (r *FieldRule) Name() lint.RuleName {
	return r.RuleName
}

// Lint accepts a Descriptor and applies LintField on it
// only if the descriptor is a FieldDescriptor.
func (r *FieldRule) Lint(d lint.Descriptor) []lint.Problem {
	if f, ok := d.(*Field); ok {
		if r.OnlyIf == nil || r.OnlyIf(f.descriptor) {
			return r.LintField(f.descriptor)
		}
	}
	return nil
}

// ServiceRule defines a lint rule that is run on each service.
type ServiceRule struct {
	RuleName lint.RuleName

	// LintService accepts a ServiceDescriptor and lints it.
	LintService func(*desc.ServiceDescriptor) []lint.Problem

	// OnlyIf accepts a ServiceDescriptor and determines whether this rule
	// is applicable.
	OnlyIf func(*desc.ServiceDescriptor) bool

	noPositional struct{}
}

// Name returns the name of the rule.
func (r *ServiceRule) Name() lint.RuleName {
	return r.RuleName
}

// Lint accepts a Descriptor and applies LintSservice on it
// only if the descriptor is a ServiceDescriptor.
func (r *ServiceRule) Lint(d lint.Descriptor) []lint.Problem {
	if s, ok := d.(*Service); ok {
		if r.OnlyIf == nil || r.OnlyIf(s.descriptor) {
			return r.LintService(s.descriptor)
		}
	}
	return nil
}

// MethodRule defines a lint rule that is run on each method.
type MethodRule struct {
	RuleName lint.RuleName

	// LintMethod accepts a MethodDescriptor and lints it.
	LintMethod func(*desc.MethodDescriptor) []lint.Problem

	// OnlyIf accepts a MethodDescriptor and determines whether this rule
	// is applicable.
	OnlyIf func(*desc.MethodDescriptor) bool

	noPositional struct{}
}

// Name returns the name of the rule.
func (r *MethodRule) Name() lint.RuleName {
	return r.RuleName
}

// Lint accepts a Descriptor and applies LintMethod on it
// only if the descriptor is a MethodDescriptor.
func (r *MethodRule) Lint(d lint.Descriptor) []lint.Problem {
	if m, ok := d.(*Method); ok {
		if r.OnlyIf == nil || r.OnlyIf(m.descriptor) {
			return r.LintMethod(m.descriptor)
		}
	}
	return nil
}

// EnumRule defines a lint rule that is run on each enum.
type EnumRule struct {
	RuleName lint.RuleName

	// LintEnum accepts a EnumDescriptor and lints it.
	LintEnum func(*desc.EnumDescriptor) []lint.Problem

	// OnlyIf accepts an EnumDescriptor and determines whether this rule
	// is applicable.
	OnlyIf func(*desc.EnumDescriptor) bool

	noPositional struct{}
}

// Name returns the name of the rule.
func (r *EnumRule) Name() lint.RuleName {
	return r.RuleName
}

// Lint accepts a Descriptor and applies LintEnum on it
// only if the descriptor is a EnumDescriptor.
func (r *EnumRule) Lint(d lint.Descriptor) []lint.Problem {
	if e, ok := d.(*Enum); ok {
		if r.OnlyIf == nil || r.OnlyIf(e.descriptor) {
			return r.LintEnum(e.descriptor)
		}
	}
	return nil
}

// EnumValueRule defines a lint rule that is run on each EnumValue.
type EnumValueRule struct {
	RuleName lint.RuleName

	// LintEnumValue accepts a EnumValueDescriptor and lints it.
	LintEnumValue func(*desc.EnumValueDescriptor) []lint.Problem

	// OnlyIf accepts an EnumValueDescriptor and determines whether this rule
	// is applicable.
	OnlyIf func(*desc.EnumValueDescriptor) bool

	noPositional struct{}
}

// Name returns the name of the rule.
func (r *EnumValueRule) Name() lint.RuleName {
	return r.RuleName
}

// Lint accepts a Descriptor and applies LintEnumValue on it
// only if the descriptor is a EnumValueDescriptor.
func (r *EnumValueRule) Lint(d lint.Descriptor) []lint.Problem {
	if e, ok := d.(*EnumValue); ok {
		if r.OnlyIf == nil || r.OnlyIf(e.descriptor) {
			return r.LintEnumValue(e.descriptor)
		}
	}
	return nil
}

// DescriptorRule defines a lint rule that is run on every descriptor
// in the file (but not the file itself).
type DescriptorRule struct {
	RuleName lint.RuleName

	// LintDescriptor accepts a generic descriptor and lints it.
	//
	// Note: Unless the descriptor is typecast to a more specific type,
	// only a subset of methods are available to it.
	LintDescriptor func(desc.Descriptor) []lint.Problem

	// OnlyIf accepts a Descriptor and determines whether this rule
	// is applicable.
	OnlyIf func(desc.Descriptor) bool

	noPositional struct{}
}

// Name returns the name of the rule.
func (r *DescriptorRule) Name() lint.RuleName {
	return r.RuleName
}

// Lint accepts a Descriptor and applies LintDescriptor on it.
func (r *DescriptorRule) Lint(d lint.Descriptor) []lint.Problem {
	if v, ok := d.(*Descriptor); ok {
		if r.OnlyIf == nil || r.OnlyIf(v.descriptor) {
			return r.LintDescriptor(v.descriptor)
		}
	}
	return nil
}
