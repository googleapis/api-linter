// Copyright 2023 Google LLC
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

package aip0121

import (
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

var noCycles = &lint.MessageRule{
	Name:   lint.NewRuleName(121, "no-cycles"),
	OnlyIf: utils.IsResource,
	LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
		var problems []lint.Problem
		references := make(map[*desc.FieldDescriptor]*desc.MessageDescriptor)

		for _, field := range m.GetFields() {
			if !isMutableReference(field) {
				continue
			}

			ref := utils.GetResourceReference(field)
			// Ignore child_type references for right now, only focus on direct
			// references.
			if ref.GetChildType() != "" {
				continue
			}
			resMsg := utils.FindResourceMessage(ref.GetType(), m.GetFile())
			// Skip unresolvable resource references, or those that don't resolve
			// to a message e.g. file-level google.api.resource_defintion.
			if resMsg == nil {
				continue
			}
			references[field] = resMsg
		}

		origRes := utils.GetResource(m)
		for origRef, refMsg := range references {
			for _, field := range refMsg.GetFields() {
				if !isMutableReference(field) {
					continue
				}
				ref := utils.GetResourceReference(field)
				if ref.GetType() == origRes.GetType() {
					problems = append(problems, lint.Problem{
						Message:    "mutable resource reference introduces a resource reference cycle",
						Descriptor: origRef,
						Location:   locations.FieldResourceReference(origRef),
					})
					break
				}
			}
		}

		return problems
	},
}

func isMutableReference(f *desc.FieldDescriptor) bool {
	behaviors := utils.GetFieldBehavior(f)
	return utils.HasResourceReference(f) && !behaviors.Contains("OUTPUT_ONLY")
}
