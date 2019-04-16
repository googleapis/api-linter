package rules

import (
	"github.com/golang/protobuf/v2/reflect/protoreflect"
	"github.com/jgeewax/api-linter/lint"
	"github.com/jgeewax/api-linter/proto"
)

type protoCheckers struct {
	CheckDescriptor          func(protoreflect.Descriptor, lint.DescriptorSource) ([]lint.Problem, error)
	CheckEnumDescriptor      func(protoreflect.EnumDescriptor, lint.DescriptorSource) ([]lint.Problem, error)
	CheckEnumValueDescriptor func(protoreflect.EnumValueDescriptor, lint.DescriptorSource) ([]lint.Problem, error)
	CheckFieldDescriptor     func(protoreflect.FieldDescriptor, lint.DescriptorSource) ([]lint.Problem, error)
	CheckFileDescriptor      func(protoreflect.FileDescriptor, lint.DescriptorSource) ([]lint.Problem, error)
	CheckMessageDescriptor   func(protoreflect.MessageDescriptor, lint.DescriptorSource) ([]lint.Problem, error)
	CheckMethodDescriptor    func(protoreflect.MethodDescriptor, lint.DescriptorSource) ([]lint.Problem, error)
	CheckServiceDescriptor   func(protoreflect.ServiceDescriptor, lint.DescriptorSource) ([]lint.Problem, error)
	CheckOneofDescriptor     func(protoreflect.OneofDescriptor, lint.DescriptorSource) ([]lint.Problem, error)

	ruleInfo ruleInfo
	problems []lint.Problem
	source   lint.DescriptorSource
}

func (c *protoCheckers) addProblems(p ...lint.Problem) {
	c.problems = append(c.problems, p...)
}

func (c *protoCheckers) ConsumeDescriptor(d protoreflect.Descriptor) error {
	if c.source.IsRuleDisabled(c.ruleInfo.Name(), d) {
		return nil
	}

	if c.CheckDescriptor != nil {
		problems, err := c.CheckDescriptor(d, c.source)
		if err != nil {
			return err
		}
		c.addProblems(problems...)
	}

	switch desc := d.(type) {
	case protoreflect.EnumDescriptor:
		if c.CheckEnumDescriptor != nil {
			problems, err := c.CheckEnumDescriptor(desc, c.source)
			if err != nil {
				return err
			}
			c.addProblems(problems...)
		}
	case protoreflect.EnumValueDescriptor:
		if c.CheckEnumValueDescriptor != nil {
			problems, err := c.CheckEnumValueDescriptor(desc, c.source)
			if err != nil {
				return err
			}
			c.addProblems(problems...)
		}
	case protoreflect.FieldDescriptor:
		if c.CheckFieldDescriptor != nil {
			problems, err := c.CheckFieldDescriptor(desc, c.source)
			if err != nil {
				return err
			}
			c.addProblems(problems...)
		}
	case protoreflect.FileDescriptor:
		if c.CheckFileDescriptor != nil {
			problems, err := c.CheckFileDescriptor(desc, c.source)
			if err != nil {
				return err
			}
			c.addProblems(problems...)
		}
	case protoreflect.MethodDescriptor:
		if c.CheckMethodDescriptor != nil {
			problems, err := c.CheckMethodDescriptor(desc, c.source)
			if err != nil {
				return err
			}
			c.addProblems(problems...)
		}
	case protoreflect.MessageDescriptor:
		if c.CheckMessageDescriptor != nil {
			problems, err := c.CheckMessageDescriptor(desc, c.source)
			if err != nil {
				return err
			}
			c.addProblems(problems...)
		}
	case protoreflect.ServiceDescriptor:
		if c.CheckServiceDescriptor != nil {
			problems, err := c.CheckServiceDescriptor(desc, c.source)
			if err != nil {
				return err
			}
			c.addProblems(problems...)
		}
	case protoreflect.OneofDescriptor:
		if c.CheckOneofDescriptor != nil {
			problems, err := c.CheckOneofDescriptor(desc, c.source)
			if err != nil {
				return err
			}
			c.addProblems(problems...)
		}
	}

	return nil
}

func (c *protoCheckers) check(req lint.Request, ri ruleInfo) (lint.Response, error) {
	f := req.ProtoFile()
	c.source = req.DescriptorSource()
	c.ruleInfo = ri
	if err := proto.WalkDescriptor(f, c); err != nil {
		return lint.Response{}, err
	}
	return lint.Response{
		Problems: c.problems,
	}, nil
}
