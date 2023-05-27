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
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
	"github.com/stoewer/go-strcase"
)

var responsePluralFirstField = &lint.MessageRule{
	Name: lint.NewRuleName(158, "response-plural-first-field"),
	OnlyIf: func(m *desc.MessageDescriptor) bool {
		return isPaginatedResponseMessage(m) && len(m.GetFields()) > 0
	},
	LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
		// Throw a linter warning if, the first field in the message is not named
		// according to plural(message_name.to_snake().split('_')[1:-1]).
		firstField := m.GetFields()[0]

		// If the field is a resource, use the `plural` annotation to decide the
		// appropriate plural field name.
		want := utils.GetResourcePlural(utils.GetResource(firstField.GetMessageType()))
		if want != "" {
			want = strcase.SnakeCase(want)
		} else {
			want = utils.ToPlural(firstField.GetName())
		}

		if want != firstField.GetName() {
			return []lint.Problem{{
				Message:    "First field of Paginated RPCs' response should be plural.",
				Suggestion: want,
				Descriptor: firstField,
				Location:   locations.DescriptorName(m),
			}}
		}

		return nil
	},
}
