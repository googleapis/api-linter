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
	"github.com/googleapis/api-linter/v2/lint"
	"github.com/googleapis/api-linter/v2/rules/internal/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var noSelfLinks = &lint.MessageRule{
	Name:   lint.NewRuleName(122, "no-self-links"),
	OnlyIf: utils.IsResource,
	LintMessage: func(m protoreflect.MessageDescriptor) []lint.Problem {
		problems := []lint.Problem{}
		for i := 0; i < m.Fields().Len(); i++ {
			field := m.Fields().Get(i)
			if field.Kind() == protoreflect.StringKind &&
				field.Name() == "self_link" {
				problems = append(problems, lint.Problem{
					Message:    "Resources must not contain self-links.",
					Descriptor: field,
				})
			}
		}
		return problems
	},
}
