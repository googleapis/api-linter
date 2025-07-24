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
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var nameNeverOptional = &lint.MessageRule{
	Name: lint.NewRuleName(123, "name-never-optional"),
	OnlyIf: func(m protoreflect.MessageDescriptor) bool {
		f := "name"
		if nf := utils.GetResource(m).GetNameField(); nf != "" {
			f = nf
		}
		return utils.IsResource(m) && m.Fields().ByName(protoreflect.Name(f)) != nil
	},
	LintMessage: func(m protoreflect.MessageDescriptor) []lint.Problem {
		f := "name"
		if nf := utils.GetResource(m).GetNameField(); nf != "" {
			f = nf
		}
		field := m.Fields().ByName(protoreflect.Name(f))

		if field.HasOptionalKeyword() {
			return []lint.Problem{{
				Message:    "Resource name fields must never be labeled with proto3_optional",
				Descriptor: field,
				Location:   locations.FieldLabel(field),
				Suggestion: "",
			}}
		}

		return nil
	},
}
