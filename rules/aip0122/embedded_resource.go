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

package aip0122

import (
	"fmt"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var embeddedResource = &lint.MessageRule{
	Name:   lint.NewRuleName(122, "embedded-resource"),
	OnlyIf: utils.IsResource,
	LintMessage: func(m protoreflect.MessageDescriptor) []lint.Problem {
		var problems []lint.Problem
		for i := 0; i < m.Fields().Len(); i++ {
			f := m.Fields().Get(i)
			mt := f.Message()
			if mt == nil || !utils.IsResource(mt) {
				continue
			}

			r := utils.GetResource(mt)
			if utils.IsResourceRevision(m) && utils.IsRevisionRelationship(r, utils.GetResource(m)) {
				continue
			}

			suggestion := fmt.Sprintf(`string %s = %d [(google.api.resource_reference).type = %q];`, f.Name(), f.Number(), r.GetType())
			if f.IsList() {
				suggestion = fmt.Sprintf("repeated %s", suggestion)
			}
			problems = append(problems, lint.Problem{
				Message:    "refer to a resource by name, not by embedding the resource message",
				Descriptor: f,
				Suggestion: suggestion,
			})
		}

		return problems
	},
}
