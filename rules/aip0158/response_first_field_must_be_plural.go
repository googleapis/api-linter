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

package aip0158

import (
	"fmt"

	"github.com/gertd/go-pluralize"
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/jhump/protoreflect/desc"
)

var responseFirstFieldMustBePlural = &lint.MessageRule{
	Name:   lint.NewRuleName(158, "response-first-field-must-be-plural"),
	OnlyIf: isPaginatedResponseMessage,
	LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
		firstField := m.FindFieldByNumber(1)
		pluralize := pluralize.NewClient()
		if !pluralize.IsPlural(firstField.GetName()) {
			want := pluralize.Plural(firstField.GetName())

			return []lint.Problem{{
				Message:    fmt.Sprintf("The first field should be plural."),
				Suggestion: want,
				Descriptor: firstField,
				Location:   locations.DescriptorName(m),
			}}
		}
		
		return nil
	},
}
