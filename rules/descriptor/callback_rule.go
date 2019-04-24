// Package descriptor contains functions that walks, consumes and checks a proto descriptor.
package descriptor

import (
	"github.com/golang/protobuf/v2/reflect/protoreflect"
	"github.com/jgeewax/api-linter/lint"
	"github.com/jgeewax/api-linter/lint/protowalk"
)

// Callback defines a callback that can be invoke on a descriptor
// with its source, which will return a list of lint.Problem or an error.
type Callback interface {
	Apply(protoreflect.Descriptor, lint.DescriptorSource) ([]lint.Problem, error)
}

// CallbackRule is a lint.Rule with a Callback that checks descriptors.
type CallbackRule struct {
	RuleInfo lint.RuleInfo
	Callback Callback

	problems []lint.Problem
	source   lint.DescriptorSource
}

// Info returns a RuleInfo for this rule.
func (r *CallbackRule) Info() lint.RuleInfo {
	return r.RuleInfo
}

// Lint accepts a lint.Request, then walks in the proto file
// by applying the Callback on each encountered descriptor,
// and finally returns a lint.Response or an error.
func (r *CallbackRule) Lint(req lint.Request) (lint.Response, error) {
	r.source = req.DescriptorSource()
	r.problems = []lint.Problem{}

	if err := protowalk.Walk(req.ProtoFile(), r); err != nil {
		return lint.Response{}, err
	}

	return lint.Response{Problems: r.problems}, nil
}

// Consume implements `Consumer` that will check the given descriptor.
func (r *CallbackRule) Consume(d protoreflect.Descriptor) error {
	problems, err := r.Callback.Apply(d, r.source)
	if err != nil {
		return err
	}
	r.addProblems(problems...)

	return nil
}

func (r *CallbackRule) addProblems(problems ...lint.Problem) {
	for _, p := range problems {
		if !p.Location.IsValid() {
			loc, _ := r.source.DescriptorLocation(p.Descriptor)
			p.Location = loc
		}
		r.problems = append(r.problems, p)
	}
}
