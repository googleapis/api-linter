package rules

import (
	"github.com/golang/protobuf/v2/reflect/protoreflect"
	"github.com/jgeewax/api-linter/lint"
	"github.com/jgeewax/api-linter/protovisit"
)

// checkers contains a collection of check functions for different descriptors and
// a set of visitors, which can be used to travel in the protobuf file.
type checkers struct {
	DescriptorCheck func(protoreflect.Descriptor, lint.DescriptorSource) []lint.Problem
	EnumCheck       func(protoreflect.EnumDescriptor, lint.DescriptorSource) []lint.Problem
	EnumValueCheck  func(protoreflect.EnumValueDescriptor, lint.DescriptorSource) []lint.Problem
	ExtensionCheck  func(protoreflect.ExtensionDescriptor, lint.DescriptorSource) []lint.Problem
	FieldCheck      func(protoreflect.FieldDescriptor, lint.DescriptorSource) []lint.Problem
	MessageCheck    func(protoreflect.MessageDescriptor, lint.DescriptorSource) []lint.Problem
	MethodCheck     func(protoreflect.MethodDescriptor, lint.DescriptorSource) []lint.Problem
	OneofCheck      func(protoreflect.OneofDescriptor, lint.DescriptorSource) []lint.Problem
	ServiceCheck    func(protoreflect.ServiceDescriptor, lint.DescriptorSource) []lint.Problem

	DescriptorVisitor protovisit.DescriptorVisitor
	EnumVisitor       protovisit.EnumVisitor
	ExtensionVisitor  protovisit.ExtensionVisitor
	MessageVisitor    protovisit.MessageVisitor
	ServiceVisitor    protovisit.ServiceVisitor

	rule     lint.Rule
	descriptorSource lint.DescriptorSource
	problems []lint.Problem
}

func (c *checkers) Lint(rule lint.Rule, req lint.Request) (lint.Response, error) {
	c.descriptorSource = req.DescriptorSource()
	f := req.ProtoFile()
	if c.DescriptorVisitor != nil {
		protovisit.WalkDescriptor(f, c.DescriptorVisitor, c)
	}
	if c.EnumVisitor != nil {
		protovisit.WalkEnum(f, c.EnumVisitor, c)
	}
	if c.ExtensionVisitor != nil {
		protovisit.WalkExtension(f, c.ExtensionVisitor, c)
	}
	if c.MessageVisitor != nil {
		protovisit.WalkMessage(f, c.MessageVisitor, c)
	}
	if c.ServiceVisitor != nil {
		protovisit.WalkService(f, c.ServiceVisitor, c)
	}
	return lint.Response{Problems: c.problems}, nil
}

func (c *checkers) VisitDescriptor(d protoreflect.Descriptor) {
	if c.DescriptorCheck != nil && c.isRuleEnabled(d) {
		c.addProblems(c.DescriptorCheck(d, c.descriptorSource)...)
	}

	switch d.(type) {
	case protoreflect.EnumDescriptor:
		c.VisitEnum(d.(protoreflect.EnumDescriptor))
	case protoreflect.EnumValueDescriptor:
		c.VisitEnumValue(d.(protoreflect.EnumValueDescriptor))
	case protoreflect.MessageDescriptor:
		c.VisitMessage(d.(protoreflect.MessageDescriptor))
	case protoreflect.MethodDescriptor:
		c.VisitMethod(d.(protoreflect.MethodDescriptor))
	case protoreflect.ServiceDescriptor:
		c.VisitService(d.(protoreflect.ServiceDescriptor))
	case protoreflect.OneofDescriptor:
		c.VisitOneof(d.(protoreflect.OneofDescriptor))
	case protoreflect.FieldDescriptor:
		f := d.(protoreflect.FieldDescriptor)
		if f.ExtendedType() != nil {
			c.VisitExtension(f)
		} else {
			c.VisitField(f)
		}
	}
}

func (c *checkers) VisitExtension(e protoreflect.ExtensionDescriptor) {
	if c.ExtensionCheck != nil && c.isRuleEnabled(e) {
		c.addProblems(c.ExtensionCheck(e, c.descriptorSource)...)
	}
}

func (c *checkers) VisitEnum(e protoreflect.EnumDescriptor) {
	if c.EnumCheck != nil && c.isRuleEnabled(e) {
		c.addProblems(c.EnumCheck(e, c.descriptorSource)...)
	}
}

func (c *checkers) VisitEnumValue(ev protoreflect.EnumValueDescriptor) {
	if c.EnumValueCheck != nil && c.isRuleEnabled(ev) {
		c.addProblems(c.EnumValueCheck(ev, c.descriptorSource)...)
	}
}

func (c *checkers) VisitField(f protoreflect.FieldDescriptor) {
	if c.FieldCheck != nil && c.isRuleEnabled(f) {
		c.addProblems(c.FieldCheck(f, c.descriptorSource)...)
	}
}

func (c *checkers) VisitMessage(m protoreflect.MessageDescriptor) {
	if c.MessageCheck != nil && c.isRuleEnabled(m) {
		c.addProblems(c.MessageCheck(m, c.descriptorSource)...)
	}
}

func (c *checkers) VisitMethod(m protoreflect.MethodDescriptor) {
	if c.MethodCheck != nil && c.isRuleEnabled(m) {
		c.addProblems(c.MethodCheck(m, c.descriptorSource)...)
	}
}

func (c *checkers) VisitOneof(o protoreflect.OneofDescriptor) {
	if c.OneofCheck != nil && c.isRuleEnabled(o) {
		c.addProblems(c.OneofCheck(o, c.descriptorSource)...)
	}
}

func (c *checkers) VisitService(s protoreflect.ServiceDescriptor) {
	if c.ServiceCheck != nil && c.isRuleEnabled(s) {
		c.addProblems(c.ServiceCheck(s, c.descriptorSource)...)
	}
}

func (c *checkers) addProblems(problems ...lint.Problem) {
	c.problems = append(c.problems, problems...)
}

func (c *checkers) isRuleEnabled(d protoreflect.Descriptor) bool {
	return !c.descriptorSource.IsRuleDisabled(c.rule.ID(), d)
}
