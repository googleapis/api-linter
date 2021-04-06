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

package aip0148

import (
	"fmt"
	"strings"

	"bitbucket.org/creachadair/stringset"
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

var declarativeFriendlyRequired = &lint.MessageRule{
	Name: lint.NewRuleName(148, "declarative-friendly-fields"),
	OnlyIf: func(m *desc.MessageDescriptor) bool {
		if resource := utils.DeclarativeFriendlyResource(m); resource == m {
			return true
		}
		return false
	},
	LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
		singleton := isSingletonResource(m)
		// Define the fields that are expected.
		missingFields := stringset.New()
		for name, typ := range reqFields {
			if singleton && singletonExceptions.Contains(name) {
				continue
			}
			f := m.FindFieldByName(name)
			if f == nil || utils.GetTypeName(f) != typ {
				missingFields.Add(fmt.Sprintf("%s %s", typ, name))
			}
		}
		if missingFields.Len() > 0 {
			msg := "Declarative-friendly resources must include the following fields:\n"
			if missingFields.Len() == 1 {
				msg = fmt.Sprintf(
					"Declarative-friendly resources must include the `%s` field.",
					missingFields.Unordered()[0],
				)
			} else {
				for _, field := range missingFields.Elements() {
					msg += fmt.Sprintf("  - `%s`\n", field)
				}
			}
			return []lint.Problem{{
				Message:    strings.TrimSuffix(msg, "\n"),
				Descriptor: m,
			}}
		}
		return nil
	},
}

var reqFields = map[string]string{
	"name":         "string",
	"uid":          "string",
	"display_name": "string",
	"create_time":  "google.protobuf.Timestamp",
	"update_time":  "google.protobuf.Timestamp",
	"delete_time":  "google.protobuf.Timestamp",
}

var singletonExceptions = stringset.New(
	"delete_time",
	"uid",
)

func isSingletonResource(m *desc.MessageDescriptor) bool {
	// If the pattern ends in something other than "}", that indicates that this is a singleton.
	//
	// For example:
	//   publishers/{publisher}/books/{book} -- not a singleton, many books
	//   publishers/*/settings -- a singleton; one settings object per publisher
	for _, pattern := range utils.GetResource(m).GetPattern() {
		if !strings.HasSuffix(pattern, "}") {
			return true
		}
	}
	return false
}
