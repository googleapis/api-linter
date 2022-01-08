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

package aip0123

import (
	"fmt"
	"strings"

	"github.com/commure/api-linter/lint"
	"github.com/commure/api-linter/locations"
	"github.com/commure/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

var resourcePattern = &lint.MessageRule{
	Name:   lint.NewRuleName(123, "resource-pattern"),
	OnlyIf: hasResourceAnnotation,
	LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
		resource := utils.GetResource(m)
		// Are any patterns declared at all? If not, complain.
		if len(resource.GetPattern()) == 0 {
			return []lint.Problem{{
				Message:    "Resources should declare resource name pattern(s).",
				Descriptor: m,
				Location:   locations.MessageResource(m),
			}}
		}

		// Ensure that the constant segments of the pattern uses camel case,
		// not snake case.
		for _, pattern := range resource.GetPattern() {
			if strings.Contains(getPlainPattern(pattern), "_") {
				return []lint.Problem{{
					Message: fmt.Sprintf(
						"Resource patterns should use camel case (apart from the variable names), such as %q.",
						getDesiredPattern(pattern),
					),
					Descriptor: m,
					Location:   locations.MessageResource(m),
				}}
			}
		}
		return nil
	},
}
