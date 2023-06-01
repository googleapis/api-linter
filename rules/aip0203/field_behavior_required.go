// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aip0203

import (
	"fmt"

	"bitbucket.org/creachadair/stringset"
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

var fbs = stringset.New(
	"OPTIONAL", "REQUIRED", "OUTPUT_ONLY", "IMMUTABLE",
)

var fieldBehaviorRequired = &lint.MethodRule{
	Name: lint.NewRuleName(203, "field-behavior-required"),
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		var ps []lint.Problem

		req := m.GetInputType()
		resp := m.GetOutputType()

		reqProblems := problems(req, 0, 0)

		var respProblems map[string][]lint.Problem
		if r := utils.GetResource(resp); r != nil {
			respProblems = problems(resp, 0, 0)
		} else {
			respProblems = problems(resp, 1, 0)
		}

		for _, p := range reqProblems {
			ps = append(ps, p...)
		}

		for msgType, p := range respProblems {
			// Avoid dups
			if _, ok := reqProblems[msgType]; !ok {
				ps = append(ps, p...)
			}
		}

		if len(ps) == 0 {
			return nil
		}

		return ps
	},
}

func problems(m *desc.MessageDescriptor, minDepth, currDepth int) map[string][]lint.Problem {
	ps := make(map[string][]lint.Problem)

	for _, f := range m.GetFields() {
		mt := f.GetMessageType()

		if minDepth <= currDepth {
			p := checkFieldBehavior(f)
			if p != nil {
				name := m.GetFullyQualifiedName()
				ps[name] = append(ps[name], *p)
			}
		}

		if mt != nil {
			for name, p := range problems(mt, minDepth, currDepth+1) {
				ps[name] = append(ps[name], p...)
			}
		}
	}

	return ps
}

func checkFieldBehavior(f *desc.FieldDescriptor) *lint.Problem {
	fb := utils.GetFieldBehavior(f)

	if len(fb) == 0 {
		return &lint.Problem{
			Message:    fmt.Sprintf("google.api.field_behavior annotation must be set on %q and contain one of, \"%v\"", f.GetName(), fbs),
			Descriptor: f,
		}
	}

	if !fbs.Intersects(fb) {
		// check for at least one valid annotation
		return &lint.Problem{
			Message: fmt.Sprintf(
				"google.api.field_behavior must contain at least one, \"%v\"", fbs),
			Descriptor: f,
		}
	}

	return nil
}
