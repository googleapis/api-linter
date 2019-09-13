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

// Package lint provides lint functions for Google APIs that register rules and user configurations,
// apply those rules to a lint request, and produce lint results.
package lint

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jhump/protoreflect/desc"
)

// Linter checks API files and returns a list of detected problems.
type Linter struct {
	rules   Rules
	configs Configs
}

// New creates and returns a linter with the given rules and configs.
func New(rules Rules, configs Configs) *Linter {
	l := &Linter{
		rules:   rules,
		configs: configs,
	}
	return l
}

// LintProtos checks protobuf files and returns a list of problems or an error.
func (l *Linter) LintProtos(files ...*desc.FileDescriptor) ([]Response, error) {
	var responses []Response
	for _, proto := range files {
		resp, err := l.lintFileDescriptor(proto)
		if err != nil {
			return nil, err
		}
		responses = append(responses, resp)
	}
	return responses, nil
}

// run executes rules on the request.
//
// It uses the proto file path to determine which rules will
// be applied to the request, according to the list of Linter
// configs.
func (l *Linter) lintFileDescriptor(fd *desc.FileDescriptor) (Response, error) {
	resp := Response{
		FilePath: fd.GetName(),
		Problems: []Problem{},
	}
	var errMessages []string

	for name, rule := range l.rules {
		var config RuleConfig

		if c, err := l.configs.GetRuleConfig(fd.GetName(), name); err == nil {
			config = config.withOverride(c)
		} else {
			errMessages = append(errMessages, err.Error())
			continue
		}

		// Run the linter rule against this file, and throw away any problems
		// which should have been disabled.
		if !config.Disabled {
			if problems, err := l.runAndRecoverFromPanics(rule, fd); err == nil {
				for _, p := range problems {
					if ruleIsEnabled(rule, p.Descriptor) {
						p.RuleID = rule.GetName()
						p.Category = config.Category
						resp.Problems = append(resp.Problems, p)
					}
				}
			} else {
				errMessages = append(errMessages, err.Error())
			}
		}
	}

	var err error
	if len(errMessages) != 0 {
		err = errors.New(strings.Join(errMessages, "; "))
	}

	return resp, err
}

func (l *Linter) runAndRecoverFromPanics(rule Rule, fd *desc.FileDescriptor) (probs []Problem, err error) {
	defer func() {
		if r := recover(); r != nil {
			if rerr, ok := r.(error); ok {
				err = rerr
			} else {
				err = fmt.Errorf("panic occurred during rule execution: %v", r)
			}
		}
	}()

	return rule.Lint(fd), nil
}
