package main

import (
	descriptorpb "github.com/golang/protobuf/v2/types/descriptor"
	"github.com/googleapis/api-linter/lint"
)

// protoParser is an interface wrapping the ParseProto method.
type protoParser interface {
	ParseProto(filenames ...string) ([]*descriptorpb.FileDescriptorProto, error)
}

// linter is an API linter.
type linter struct {
	rt *lint.Runtime
}

// newLinter creates and returns a linter with the given rules and configs.
func newLinter(rules []lint.Rule, configs lint.RuntimeConfigs) *linter {
	l := &linter{
		rt: lint.NewRuntime(),
	}
	l.rt.AddRules(rules...)
	l.rt.AddConfigs(configs...)
	return l
}

// LintProto checks protobuf files and returns a list of problems or an error.
func (l *linter) LintProto(files []*descriptorpb.FileDescriptorProto) ([]lint.Problem, error) {
	return checkProto(l.rt, files)
}

func checkProto(rt *lint.Runtime, files []*descriptorpb.FileDescriptorProto) ([]lint.Problem, error) {
	var problems []lint.Problem
	for _, proto := range files {
		req, err := lint.NewProtoRequest(proto)
		if err != nil {
			return nil, err
		}
		resp, err := rt.Run(req)
		if err != nil {
			return nil, err
		}
		for _, prob := range resp.Problems {
			problems = append(problems, prob)
		}
	}
	return problems, nil
}
