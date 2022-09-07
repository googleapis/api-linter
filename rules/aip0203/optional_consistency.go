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

package aip0203

import (
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

// If a message has a field which is described as optional, ensure that every
// optional field on the message has this indicator. Oneof fields do not count.
var optionalBehaviorConsistency = &lint.MessageRule{
	Name:   lint.NewRuleName(203, "optional-consistency"),
	OnlyIf: messageHasOptionalFieldBehavior,
	LintMessage: func(m *desc.MessageDescriptor) (problems []lint.Problem) {
		for _, f := range m.GetFields() {
			if utils.GetFieldBehavior(f).Len() == 0 && !standardFields.Contains(f.GetName()) && f.GetOneOf() == nil {
				problems = append(problems, lint.Problem{
					Message:    "Within a single message, either all optional fields should be indicated, or none of them should be.",
					Descriptor: f,
				})
			}
		}
		return
	},
}
