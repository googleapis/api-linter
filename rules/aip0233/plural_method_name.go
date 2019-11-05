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

package aip0233

import (
	"fmt"

	"github.com/gertd/go-pluralize"
	"github.com/googleapis/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
)

var pluralMethodName = &lint.MethodRule{
	Name:   lint.NewRuleName("core", "0233", "plural-method-name"),
	OnlyIf: isBatchCreateMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		// Note: m.GetName()[11:] is used to retrieve the resource name from the
		// method name. For example, "BatchCreateFoos" -> "Foos"
		pluralMethodResourceName := m.GetName()[11:]

		pluralize := pluralize.NewClient()

		// Rule check: Establish that for methods such as `BatchCreateFoos`
		if !pluralize.IsPlural(pluralMethodResourceName) {
			return []lint.Problem{{
				Message: fmt.Sprintf(
					`The resource part in method name %q shouldn't be %q, but should be its plural form %q`,
					m.GetName(), pluralMethodResourceName, pluralize.Plural(pluralMethodResourceName),
				),
				Descriptor: m,
			}}
		}

		return nil
	},
}
