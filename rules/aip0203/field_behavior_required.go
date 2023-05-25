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

var minimumRequiredFieldBehavior = stringset.New(
	"OPTIONAL", "REQUIRED", "OUTPUT_ONLY", "IMMUTABLE",
)

var fieldBehaviorRequired = &lint.MethodRule{
	Name: lint.NewRuleName(203, "field-behavior-required"),
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		// we only check requests, as OutputTypes are always
		// OUTPUT_ONLY
		it := m.GetInputType()
		return checkFields(it)
	},
}

func checkFields(m *desc.MessageDescriptor) []lint.Problem {
	problems := []lint.Problem{}
	for _, f := range m.GetFields() {
		asMessage := f.GetMessageType()
		if asMessage != nil {
			problems = append(problems, checkFields(asMessage)...)
		}
		problems = append(problems, checkFieldBehavior(f)...)
	}
	return problems
}

func checkFieldBehavior(f *desc.FieldDescriptor) []lint.Problem {
	problems := []lint.Problem{}
	fieldBehavior := utils.GetFieldBehavior(f)
	if len(fieldBehavior) == 0 {
		problems = append(problems, lint.Problem{
			Message:    fmt.Sprintf("google.api.field_behavior annotation must be set, and have one of %v", minimumRequiredFieldBehavior),
			Descriptor: f,
		})
		// check for at least one valid annotation
	} else if !minimumRequiredFieldBehavior.Intersects(fieldBehavior) {
		problems = append(problems, lint.Problem{
			Message: fmt.Sprintf(
				"google.api.field_behavior must have at least one of the following behaviors set: %v", minimumRequiredFieldBehavior),
			Descriptor: f,
		})
	}
	return problems
}
