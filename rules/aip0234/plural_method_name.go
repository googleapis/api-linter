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

package aip0234

import (
	"fmt"
	"strings"

	"github.com/commure/api-linter/lint"
	"github.com/commure/api-linter/locations"
	"github.com/gertd/go-pluralize"
	"github.com/jhump/protoreflect/desc"
)

var pluralMethodName = &lint.MethodRule{
	Name:   lint.NewRuleName(234, "plural-method-name"),
	OnlyIf: isBatchUpdateMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		// Note: Retrieve the resource name from the method name. For example,
		// "BatchUpdateFoos" -> "Foos"
		pluralMethodResourceName := strings.TrimPrefix(m.GetName(), "BatchUpdate")

		pluralize := pluralize.NewClient()

		// Rule check: Establish that for methods such as `BatchUpdateFoos`
		if !pluralize.IsPlural(pluralMethodResourceName) {
			return []lint.Problem{{
				Message: fmt.Sprintf(
					`The resource part in method name %q shouldn't be %q, but should be its plural form %q`,
					m.GetName(), pluralMethodResourceName, pluralize.Plural(pluralMethodResourceName),
				),
				Descriptor: m,
				Location:   locations.DescriptorName(m),
				Suggestion: "BatchUpdate" + pluralize.Plural(pluralMethodResourceName),
			}}
		}

		return nil
	},
}
