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
	"strings"

	"bitbucket.org/creachadair/stringset"
	"github.com/googleapis/api-linter/v2/lint"
	"github.com/googleapis/api-linter/v2/locations"
	"github.com/googleapis/api-linter/v2/rules/internal/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var noMutableCycles = &lint.MessageRule{
	Name:   lint.NewRuleName(121, "no-mutable-cycles"),
	OnlyIf: utils.IsResource,
	LintMessage: func(m protoreflect.MessageDescriptor) []lint.Problem {
		res := utils.GetResource(m)

		return findCycles(res.GetType(), m, stringset.New(), nil)
	},
}

func findCycles(start string, node protoreflect.MessageDescriptor, seen stringset.Set, chain []string) []lint.Problem {
	var problems []lint.Problem
	nodeRes := utils.GetResource(node)

	chain = append(chain, nodeRes.GetType())
	seen.Add(nodeRes.GetType())

	for i := 0; i < node.Fields().Len(); i++ {
		f := node.Fields().Get(i)
		if !isMutableReference(f) {
			continue
		}
		ref := utils.GetResourceReference(f)
		// Skip indirect references for now.
		if ref.GetChildType() != "" {
			continue
		}
		if ref.GetType() == start {
			cycle := strings.Join(append(chain, start), " > ")
			problems = append(problems, lint.Problem{
				Message:    "mutable resource reference introduces a reference cycle:\n" + cycle,
				Descriptor: f,
				Location:   locations.FieldResourceReference(f),
			})
		} else if !seen.Contains(ref.GetType()) {
			next := utils.FindResourceMessage(ref.GetType(), node.ParentFile())
			// Skip unresolvable references.
			if next == nil {
				continue
			}
			if probs := findCycles(start, next, seen.Clone(), chain); len(probs) == 1 {
				// A recursive call will only have one finding as it returns
				// immediately.
				p := probs[0]
				p.Descriptor = f
				p.Location = locations.FieldResourceReference(f)

				problems = append(problems, p)
			}
		}
		// Recursive calls should return immediately upon finding a cycle.
		// The initial scan of the starting resource message should scan
		// all of its own fields and not return immediately.
		if len(problems) > 0 && len(chain) > 1 {
			return problems
		}
	}

	return problems
}

func isMutableReference(f protoreflect.FieldDescriptor) bool {
	behaviors := utils.GetFieldBehavior(f)
	return utils.HasResourceReference(f) && !behaviors.Contains("OUTPUT_ONLY")
}
