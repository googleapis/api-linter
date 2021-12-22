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

package aip0135

import (
	"fmt"

	"github.com/googleapis/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
)

// Delete methods should not have unrecognized fields.
var unknownFields = &lint.MessageRule{
	Name:   lint.NewRuleName(135, "request-unknown-fields"),
	OnlyIf: isDeleteRequestMessage,
	LintMessage: func(m *desc.MessageDescriptor) (problems []lint.Problem) {
		// Rule check: Establish that there are no unexpected fields.
		allowedFields := map[string]struct{}{
			"resource_name":          {}, // AIP-135
			"force":         {}, // AIP-135
			"allow_missing": {}, // AIP-135
			"etag":          {}, // AIP-154
			"request_id":    {}, // AIP-155
			"validate_only": {}, // AIP-163
		}
		for _, field := range m.GetFields() {
			if _, ok := allowedFields[string(field.GetName())]; !ok {
				problems = append(problems, lint.Problem{
					Message: fmt.Sprintf(
						"Unexpected field: Delete RPCs must only contain fields explicitly described in AIPs, not %q.",
						string(field.GetName()),
					),
					Descriptor: field,
				})
			}
		}

		return
	},
}
