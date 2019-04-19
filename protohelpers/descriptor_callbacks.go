package protohelpers

import (
	"github.com/golang/protobuf/v2/reflect/protoreflect"
	"github.com/jgeewax/api-linter/lint"
)

// DescriptorCallbacks implements both `Rule` and `DescriptorConsumer`. Set any of the Callback fields
// to a non nil value, and that function will be called for every encountered instance of the relevant
// descriptor.
type DescriptorCallbacks struct {
	RuleInfo lint.RuleInfo

	DescriptorCallback          func(protoreflect.Descriptor, lint.DescriptorSource) ([]lint.Problem, error)
	EnumDescriptorCallback      func(protoreflect.EnumDescriptor, lint.DescriptorSource) ([]lint.Problem, error)
	EnumValueCallback           func(protoreflect.EnumValueDescriptor, lint.DescriptorSource) ([]lint.Problem, error)
	FieldDescriptorCallback     func(protoreflect.FieldDescriptor, lint.DescriptorSource) ([]lint.Problem, error)
	ExtensionDescriptorCallback func(protoreflect.ExtensionDescriptor, lint.DescriptorSource) ([]lint.Problem, error)
	FileDescriptorCallback      func(protoreflect.FileDescriptor, lint.DescriptorSource) ([]lint.Problem, error)
	MessageDescriptorCallback   func(protoreflect.MessageDescriptor, lint.DescriptorSource) ([]lint.Problem, error)
	MethodDescriptorCallback    func(protoreflect.MethodDescriptor, lint.DescriptorSource) ([]lint.Problem, error)
	ServiceDescriptorCallback   func(protoreflect.ServiceDescriptor, lint.DescriptorSource) ([]lint.Problem, error)
	OneofDescriptorCallback     func(protoreflect.OneofDescriptor, lint.DescriptorSource) ([]lint.Problem, error)

	problems []lint.Problem
	source   lint.DescriptorSource
}

// Info implements `lint.Rule`
func (c *DescriptorCallbacks) Info() lint.RuleInfo {
	return c.RuleInfo
}

// Lint implements `lint.Rule`
func (c *DescriptorCallbacks) Lint(req lint.Request) (lint.Response, error) {
	f := req.ProtoFile()

	c.source = req.DescriptorSource()

	if err := WalkDescriptor(f, c); err != nil {
		return lint.Response{}, err
	}

	return lint.Response{
		Problems: c.problems,
	}, nil
}

// ConsumeDescriptor implements `DescriptorConsumer`
func (c *DescriptorCallbacks) ConsumeDescriptor(d protoreflect.Descriptor) error {
	if c.source.IsRuleDisabled(c.Info().Name, d) {
		return nil
	}

	if c.DescriptorCallback != nil {
		problems, err := c.DescriptorCallback(d, c.source)
		if err != nil {
			return err
		}
		c.addProblems(problems...)
	}

	switch desc := d.(type) {
	case protoreflect.EnumDescriptor:
		if c.EnumDescriptorCallback != nil {
			problems, err := c.EnumDescriptorCallback(desc, c.source)
			if err != nil {
				return err
			}
			c.addProblems(problems...)
		}
	case protoreflect.EnumValueDescriptor:
		if c.EnumValueCallback != nil {
			problems, err := c.EnumValueCallback(desc, c.source)
			if err != nil {
				return err
			}
			c.addProblems(problems...)
		}
	case protoreflect.FieldDescriptor:
		if desc.ExtendedType() != nil {
			if c.ExtensionDescriptorCallback != nil {
				problems, err := c.ExtensionDescriptorCallback(desc, c.source)
				if err != nil {
					return err
				}
				c.addProblems(problems...)
			}
		} else if c.FieldDescriptorCallback != nil {
			problems, err := c.FieldDescriptorCallback(desc, c.source)
			if err != nil {
				return err
			}
			c.addProblems(problems...)
		}
	case protoreflect.FileDescriptor:
		if c.FileDescriptorCallback != nil {
			problems, err := c.FileDescriptorCallback(desc, c.source)
			if err != nil {
				return err
			}
			c.addProblems(problems...)
		}
	case protoreflect.MethodDescriptor:
		if c.MethodDescriptorCallback != nil {
			problems, err := c.MethodDescriptorCallback(desc, c.source)
			if err != nil {
				return err
			}
			c.addProblems(problems...)
		}
	case protoreflect.MessageDescriptor:
		if c.MessageDescriptorCallback != nil {
			problems, err := c.MessageDescriptorCallback(desc, c.source)
			if err != nil {
				return err
			}
			c.addProblems(problems...)
		}
	case protoreflect.ServiceDescriptor:
		if c.ServiceDescriptorCallback != nil {
			problems, err := c.ServiceDescriptorCallback(desc, c.source)
			if err != nil {
				return err
			}
			c.addProblems(problems...)
		}
	case protoreflect.OneofDescriptor:
		if c.OneofDescriptorCallback != nil {
			problems, err := c.OneofDescriptorCallback(desc, c.source)
			if err != nil {
				return err
			}
			c.addProblems(problems...)
		}
	}

	return nil
}

func (c *DescriptorCallbacks) addProblems(p ...lint.Problem) {
	c.problems = append(c.problems, p...)
}
