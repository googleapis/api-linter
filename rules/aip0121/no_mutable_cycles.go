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
	"bitbucket.org/creachadair/stringset"
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
	rpb "google.golang.org/genproto/googleapis/api/annotations"
)

var noMutableCycles = &lint.MessageRule{
	Name:   lint.NewRuleName(121, "no-mutable-cycles"),
	OnlyIf: utils.IsResource,
	LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
		var problems []lint.Problem
		res := utils.GetResource(m)

		for _, field := range m.GetFields() {
			if !isMutableReference(field) {
				continue
			}

			ref := utils.GetResourceReference(field)
			// Focus on direct type references.
			if ref.GetChildType() != "" {
				continue
			}

			if p := findCycle(res, m, ref.GetType(), stringset.New()); p != nil {
				p.Descriptor = field
				p.Location = locations.FieldResourceReference(field)
				problems = append(problems, *p)
			}
		}

		return problems
	},
}

func findCycle(orig *rpb.ResourceDescriptor, node *desc.MessageDescriptor, nodeRef string, seen stringset.Set) *lint.Problem {
	refMsg := utils.FindResourceMessage(nodeRef, node.GetFile())
	if refMsg == nil {
		// Unable to resolve the resource reference to a message.
		return nil
	}
	seen.Add(nodeRef)

	for _, f := range refMsg.GetFields() {
		if !isMutableReference(f) {
			continue
		}
		ref := utils.GetResourceReference(f)
		// Focus on direct type references.
		if ref.GetChildType() != "" {
			continue
		}
		if ref.GetType() == orig.GetType() {
			return &lint.Problem{
				Message: "mutable resource reference introduces a resource reference cycle",
			}
		}
		if !seen.Contains(ref.GetType()) {
			if p := findCycle(orig, refMsg, ref.GetType(), seen); p != nil {
				return p
			}
		}
	}

	return nil
}

func isMutableReference(f *desc.FieldDescriptor) bool {
	behaviors := utils.GetFieldBehavior(f)
	return utils.HasResourceReference(f) && !behaviors.Contains("OUTPUT_ONLY")
}
