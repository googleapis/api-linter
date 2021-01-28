// Copyright 2021 Google LLC
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

package aip0162

import (
	"fmt"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"
)

// The Tag Revision request message should have a tag field.
var tagRevisionRequestTagField = &lint.MessageRule{
	Name:   lint.NewRuleName(162, "tag-revision-request-tag-field"),
	OnlyIf: isTagRevisionRequestMessage,
	LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
		// Rule check: Establish that a `tag` field is present.
		tagField := m.FindFieldByName("tag")
		if tagField == nil {
			return []lint.Problem{{
				Message:    fmt.Sprintf("Message %q has no `tag` field.", m.GetName()),
				Descriptor: m,
			}}
		}

		// Rule check: Establish that the tag field is a string.
		if tagField.GetType() != builder.FieldTypeString().GetType() || tagField.IsRepeated() {
			return []lint.Problem{{
				Message:    "`tag` field on Tag Revision request message must be a singular string.",
				Descriptor: tagField,
				Location:   locations.FieldType(tagField),
				Suggestion: "string",
			}}
		}

		return nil
	},
}