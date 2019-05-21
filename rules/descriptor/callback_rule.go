// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package descriptor provides helper utilities for linting proto files.
package descriptor

import (
	"google.golang.org/protobuf/reflect/protoreflect"
	"github.com/googleapis/api-linter/lint"
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
// and finally returns a list of problems or an error.
func (r *CallbackRule) Lint(req lint.Request) ([]lint.Problem, error) {
	r.source = req.DescriptorSource()
	r.problems = []lint.Problem{}

	if err := Walk(req.ProtoFile(), r); err != nil {
		return nil, err
	}

	return r.problems, nil
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
