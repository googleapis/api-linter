package rules

import (
	"github.com/golang/protobuf/v2/reflect/protoreflect"
	"github.com/jgeewax/api-linter/lint"
	"github.com/jgeewax/api-linter/proto"
)

// linter defines an operation that checks a lint Request with a Rule,
// and returns a Response and if applicable, an error.
type linter interface {
	Lint(lint.Request, lint.Rule) (lint.Response, error)
}

type protoLinter struct {
	check    descCheckFunc
	info     ruleInfo
	problems []lint.Problem
	source   lint.DescriptorSource
}

func (l *protoLinter) addProblems(p ...lint.Problem) {
	l.problems = append(l.problems, p...)
}

func (l *protoLinter) ConsumeDescriptor(d protoreflect.Descriptor) error {
	if l.source.IsRuleDisabled(l.info.Name, d) {
		return nil
	}

	problems, err := l.check(d, l.source)
	if err != nil {
		return err
	}
	l.addProblems(problems...)

	return nil
}

func (l *protoLinter) Lint(req lint.Request, rule lint.Rule) (lint.Response, error) {
	f := req.ProtoFile()
	l.source = req.DescriptorSource()
	if err := proto.WalkDescriptor(f, l); err != nil {
		return lint.Response{}, nil
	}
	return lint.Response{
		Problems: l.problems,
	}, nil
}

func newProtoLinter(info ruleInfo, check descCheckFunc) linter {
	return &protoLinter{
		check: check,
		info:  info,
	}
}
