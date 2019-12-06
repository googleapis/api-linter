// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aip0133

import (
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

var requestParentBehavior = &lint.FieldRule{
	Name: lint.NewRuleName(133, "request-parent-behavior"),
	OnlyIf: func(f *desc.FieldDescriptor) bool {
		return isCreateRequestMessage(f.GetOwner()) && f.GetName() == "parent"
	},
	LintField: func(f *desc.FieldDescriptor) []lint.Problem {
		if !utils.GetFieldBehavior(f).Contains("REQUIRED") {
			return []lint.Problem{{
				Message:    "Create requests: The `parent` field should include `(google.api.field_behavior) = REQUIRED`.",
				Descriptor: f,
			}}
		}
		return nil
	},
}
