// Copyright 2022 Google LLC
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
	"regexp"
	"unicode"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

var (
	resourceNameRegex = regexp.MustCompile(`^[\p{L}A-Za-z0-9]+$`)
)

var resourceTypeName = &lint.MessageRule{
	Name: lint.NewRuleName(123, "resource-type-name"),
	OnlyIf: func(m *desc.MessageDescriptor) bool {
		return utils.GetResource(m) != nil
	},
	LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
		resource := utils.GetResource(m)
		_, typeName, ok := utils.SplitResourceTypeName(resource.GetType())
		if !ok {
			return []lint.Problem{{
				Message:    "Resource type names must be of the form {Service Name}/{Type}.",
				Descriptor: m,
				Location:   locations.MessageResource(m),
			}}
		}
		if !unicode.IsUpper(rune(typeName[0])) {
			return []lint.Problem{{
				Message:    fmt.Sprintf("Type %q must be UpperCamelCase", typeName),
				Descriptor: m,
				Location:   locations.MessageResource(m),
			}}
		}
		if !resourceNameRegex.MatchString(typeName) {
			return []lint.Problem{{
				Message: fmt.Sprintf(
					"Type %q must only contain alphanumeric characters", typeName,
				),
				Descriptor: m,
				Location:   locations.MessageResource(m),
			}}
		}
		return nil
	},
}
