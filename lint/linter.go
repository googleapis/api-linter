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

	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
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

// LintProtos checks protobuf files and returns a list of problems or an error. Note that a *FileDescriptorProto
// for any imported file must be present in files. If any file in files has an import that is not also in the
// slice, an error will be returned.
func (l *Linter) LintProtos(files []*descriptorpb.FileDescriptorProto) ([]Response, error) {
	reg, err := makeRegistryFromAllFiles(files)

	if err != nil {
		return nil, err
	}

	return l.LintProtosWithRegistry(files, reg)
}

// LintProtosWithRegistry checks protobuf files and returns a list of Responses or returns an error.
// The input *protoregistry.Files must include all imports from any file in files.
func (l *Linter) LintProtosWithRegistry(files []*descriptorpb.FileDescriptorProto, reg *protoregistry.Files) ([]Response, error) {
	var responses []Response
	for _, proto := range files {
		req, err := NewProtoRequest(proto, reg)
		if err != nil {
			return nil, err
		}
		resp, err := l.run(req)
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
func (l *Linter) run(req Request) (Response, error) {
	resp := Response{
		FilePath: req.ProtoFile().Path(),
	}
	var errMessages []string

	for name, rl := range l.rules {
		var config RuleConfig

		if c, err := l.configs.GetRuleConfig(req.ProtoFile().Path(), name); err == nil {
			config = config.withOverride(c)
		} else {
			errMessages = append(errMessages, err.Error())
			continue
		}

		if !config.Disabled && !req.DescriptorSource().isRuleDisabledInFile(rl.Info().Name) {
			if problems, err := l.runAndRecoverFromPanics(rl, req); err == nil {
				for _, p := range problems {
					if !req.DescriptorSource().isRuleDisabled(rl.Info().Name, p.Descriptor) {
						p.RuleID = rl.Info().Name
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

func (l *Linter) runAndRecoverFromPanics(rl Rule, req Request) (probs []Problem, err error) {
	defer func() {
		if r := recover(); r != nil {
			if rerr, ok := r.(error); ok {
				err = rerr
			} else {
				err = fmt.Errorf("panic occurred during rule execution: %v", r)
			}
		}
	}()

	return rl.Lint(req)
}
