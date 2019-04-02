package rules

import (
	"github.com/golang/protobuf/v2/reflect/protoreflect"
	"github.com/jgeewax/api-linter/lint"
	"github.com/jgeewax/api-linter/protovisit"
)

// checkers contains a collection of check functions for different descriptors and
// a set of visitors, which can be used to travel in the protobuf file.
type checkers struct {
	DescriptorCheck func(protoreflect.Descriptor, lint.Context) []lint.Problem
	EnumCheck       func(protoreflect.EnumDescriptor, lint.Context) []lint.Problem
	EnumValueCheck  func(protoreflect.EnumValueDescriptor, lint.Context) []lint.Problem
	ExtensionCheck  func(protoreflect.ExtensionDescriptor, lint.Context) []lint.Problem
	FieldCheck      func(protoreflect.FieldDescriptor, lint.Context) []lint.Problem
	MessageCheck    func(protoreflect.MessageDescriptor, lint.Context) []lint.Problem
	MethodCheck     func(protoreflect.MethodDescriptor, lint.Context) []lint.Problem
	OneofCheck      func(protoreflect.OneofDescriptor, lint.Context) []lint.Problem
	ServiceCheck    func(protoreflect.ServiceDescriptor, lint.Context) []lint.Problem

	DescriptorVisitor protovisit.DescriptorVisitor
	EnumVisitor       protovisit.EnumVisitor
	ExtensionVisitor  protovisit.ExtensionVisitor
	MessageVisitor    protovisit.MessageVisitor
	ServiceVisitor    protovisit.ServiceVisitor

	rule     lint.Rule
	ctx      lint.Context
	problems []lint.Problem
}

func (c *checkers) Lint(rule lint.Rule, req lint.Request) (lint.Response, error) {
	c.ctx = req.Context()
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
		c.addProblems(c.DescriptorCheck(d, c.ctx)...)
	}
}

func (c *checkers) VisitExtension(e protoreflect.ExtensionDescriptor) {
	if c.ExtensionCheck != nil && c.isRuleEnabled(e) {
		c.addProblems(c.ExtensionCheck(e, c.ctx)...)
	}
}

func (c *checkers) VisitEnum(e protoreflect.EnumDescriptor) {
	if c.EnumCheck != nil && c.isRuleEnabled(e) {
		c.addProblems(c.EnumCheck(e, c.ctx)...)
	}
}

func (c *checkers) VisitEnumValue(ev protoreflect.EnumValueDescriptor) {
	if c.EnumValueCheck != nil && c.isRuleEnabled(ev) {
		c.addProblems(c.EnumValueCheck(ev, c.ctx)...)
	}
}

func (c *checkers) VisitField(f protoreflect.FieldDescriptor) {
	if c.FieldCheck != nil && c.isRuleEnabled(f) {
		c.addProblems(c.FieldCheck(f, c.ctx)...)
	}
}

func (c *checkers) VisitMessage(m protoreflect.MessageDescriptor) {
	if c.MessageCheck != nil && c.isRuleEnabled(m) {
		c.addProblems(c.MessageCheck(m, c.ctx)...)
	}
}

func (c *checkers) VisitMethod(m protoreflect.MethodDescriptor) {
	if c.MethodCheck != nil && c.isRuleEnabled(m) {
		c.addProblems(c.MethodCheck(m, c.ctx)...)
	}
}

func (c *checkers) VisitOneof(o protoreflect.OneofDescriptor) {
	if c.OneofCheck != nil && c.isRuleEnabled(o) {
		c.addProblems(c.OneofCheck(o, c.ctx)...)
	}
}

func (c *checkers) VisitService(s protoreflect.ServiceDescriptor) {
	if c.ServiceCheck != nil && c.isRuleEnabled(s) {
		c.addProblems(c.ServiceCheck(s, c.ctx)...)
	}
}

func (c *checkers) addProblems(problems ...lint.Problem) {
	c.problems = append(c.problems, problems...)
}

func (c *checkers) isRuleEnabled(d protoreflect.Descriptor) bool {
	return !c.ctx.DescriptorSource().IsRuleDisabled(c.rule.ID(), d)
}
