package rules

import (
	"github.com/golang/protobuf/v2/reflect/protoreflect"
	"github.com/jgeewax/api-linter/lint"
	"github.com/jgeewax/api-linter/proto"
)

// Linter defines an operation that checks a lint Request with a Rule,
// and returns a Response and if applicable, an error.
type Linter interface {
	Lint(lint.Request, lint.Rule) (lint.Response, error)
}

type protoLinter struct {
	checker  DescriptorChecker
	info     RuleInfo
	problems []lint.Problem
	source   lint.DescriptorSource
}

func (l *protoLinter) addProblems(p ...lint.Problem) {
	l.problems = append(l.problems, p...)
}

func (l *protoLinter) Consume(d protoreflect.Descriptor) error {
	if l.source.IsRuleDisabled(l.info.Name, d) {
		return nil
	}

	problems, err := l.checker.Check(d)
	if err != nil {
		return err
	}
	l.addProblems(problems...)

	return nil
}

func (l *protoLinter) Lint(req lint.Request, rule lint.Rule) (lint.Response, error) {
	f := req.ProtoFile()
	if err := proto.Walk(f, l); err != nil {
		return lint.Response{}, nil
	}
	return lint.Response{
		Problems: l.problems,
	}, nil
}

func newProtoLinter(info RuleInfo, checker DescriptorChecker) Linter {
	return &protoLinter{
		checker: checker,
		info:    info,
	}
}
