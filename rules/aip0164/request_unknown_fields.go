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

package aip0164

import (
	"fmt"

	"github.com/googleapis/api-linter/lint"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Undelete methods should not have unrecognized fields.
var requestUnknownFields = &lint.MessageRule{
	Name:   lint.NewRuleName(164, "request-unknown-fields"),
	OnlyIf: isUndeleteRequestMessage,
	LintMessage: func(m protoreflect.MessageDescriptor) (problems []lint.Problem) {
		// Rule check: Establish that there are no unexpected fields.
		allowedFields := map[string]struct{}{
			"name":          {}, // AIP-164
			"etag":          {}, // AIP-154
			"request_id":    {}, // AIP-155
			"validate_only": {}, // AIP-163
		}
		for _, field := range m.Fields() {
			if _, ok := allowedFields[string(field.Name())]; !ok {
				problems = append(problems, lint.Problem{
					Message: fmt.Sprintf(
						"Unexpected field: Undelete requests must only contain fields explicitly described in AIPs, not %q.",
						string(field.Name()),
					),
					Descriptor: field,
				})
			}
		}

		return
	},
}
