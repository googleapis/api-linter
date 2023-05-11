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
)

var resourceSingular = &lint.MessageRule{
	Name:   lint.NewRuleName(123, "resource-singular"),
	OnlyIf: hasResourceAnnotation,
	LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
		r := utils.GetResource(m)
		l := locations.MessageResource(m)
		s := r.GetSingular()
		_, typeName, ok := utils.SplitResourceTypeName(r.GetType())
		lowerTypeName := utils.ToLowerCamelCase(typeName)
		if s == "" {
			return []lint.Problem{{
				Message:    fmt.Sprintf("Resources should declare singular: %q", lowerTypeName),
				Descriptor: m,
				Location:   l,
			}}
		}
		if !ok || lowerTypeName != s {
			return []lint.Problem{{
				Message: fmt.Sprintf(
					"Resource singular should be lower camel case of type: %q",
					lowerTypeName,
				),
				Descriptor: m,
				Location:   l,
			}}
		}
		return nil
	},
}
