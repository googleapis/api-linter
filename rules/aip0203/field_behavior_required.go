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
	"OPTIONAL", "REQUIRED", "OUTPUT_ONLY",
)

var fieldBehaviorRequired = &lint.FieldRule{
	Name: lint.NewRuleName(203, "field-behavior-required"),
	LintField: func(f *desc.FieldDescriptor) []lint.Problem {
		fieldBehavior := utils.GetFieldBehavior(f)
		if len(fieldBehavior) == 0 {
			return []lint.Problem{{
				Message:    "google.api.field_behavior annotation must be set",
				Descriptor: f,
			}}
		}
		// check for at least one valid annotation
		if !minimumRequiredFieldBehavior.Intersects(fieldBehavior) {
			return []lint.Problem{{
				Message: fmt.Sprintf(
					"google.api.field_behavior must have at least one of the following behaviors set: %v",
					minimumRequiredFieldBehavior,
				),
				Descriptor: f,
			}}
		}
		return nil
	},
}
