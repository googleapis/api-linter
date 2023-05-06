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

package aip0123

import (
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

var resourcePlural = &lint.MessageRule{
	Name:   lint.NewRuleName(123, "resource-plural"),
	OnlyIf: hasResourceAnnotation,
	LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
		r := utils.GetResource(m)
		l := locations.MessageResource(m)
		p := r.GetPlural()
		if p == "" {
			return []lint.Problem{{
				Message:    "Resources should declare plural.",
				Descriptor: m,
				Location:   l,
			}}
		}
		return nil
	},
}
