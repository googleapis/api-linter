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
	"fmt"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
	"github.com/stoewer/go-strcase"
)

var resourceSingular = &lint.MessageRule{
	Name: lint.NewRuleName(123, "resource-singular"),
	OnlyIf: func(m *desc.MessageDescriptor) bool {
		r := utils.GetResource(m)
		return r != nil || r.GetType() == ""
	},
	LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
		r := utils.GetResource(m)
		l := locations.MessageResource(m)
		s := r.GetSingular()
		if s == "" {
			return []lint.Problem{{
				Message:    "Resources should declare singular.",
				Descriptor: m,
				Location:   l,
			}}
		}
		_, typeName, ok := utils.SplitResourceTypeName(r.GetType())
		lowerTypeName := strcase.LowerCamelCase(typeName)
		if !ok || lowerTypeName != s {
			return []lint.Problem{{
				Message: fmt.Sprintf(
					"Resource singular should be lower camel case of type: %q, not %q.",
					lowerTypeName, s,
				),
				Descriptor: m,
				Location:   l,
			}}
		}
		return nil
	},
}
