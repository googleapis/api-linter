// Copyright 2020 Google LLC
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

package aip0235

import (
	"fmt"
	"strings"

	"github.com/gertd/go-pluralize"
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var pluralMethodName = &lint.MethodRule{
	Name:   lint.NewRuleName(235, "plural-method-name"),
	OnlyIf: isBatchDeleteMethod,
	LintMethod: func(m protoreflect.MethodDescriptor) []lint.Problem {
		// Note: Retrieve the resource name from the method name. For example,
		// "BatchDeleteFoos" -> "Foos"
		pluralMethodResourceName := strings.TrimPrefix(string(m.Name()), "BatchDelete")

		pluralize := pluralize.NewClient()

		// Rule check: Establish that for methods such as `BatchDeleteFoos`
		if !pluralize.IsPlural(pluralMethodResourceName) {
			return []lint.Problem{{
				Message: fmt.Sprintf(
					`The resource part in method name %q shouldn't be %q, but should be its plural form %q`,
					m.Name(), pluralMethodResourceName, pluralize.Plural(pluralMethodResourceName),
				),
				Descriptor: m,
				Location:   locations.DescriptorName(m),
				Suggestion: "BatchDelete" + pluralize.Plural(pluralMethodResourceName),
			}}
		}

		return nil
	},
}
