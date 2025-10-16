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

package aip0133

import (
	"fmt"
	"strings"

	"bitbucket.org/creachadair/stringset"
	"github.com/googleapis/api-linter/v2/lint"
	"github.com/googleapis/api-linter/v2/rules/internal/utils"
	"github.com/stoewer/go-strcase"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// The create request message should not have unrecognized fields.
var unknownFields = &lint.MessageRule{
	Name:   lint.NewRuleName(133, "request-unknown-fields"),
	OnlyIf: utils.IsCreateRequestMessage,
	LintMessage: func(m protoreflect.MessageDescriptor) (problems []lint.Problem) {
		resourceMsgName := getResourceMsgNameFromReq(m)

		// Rule check: Establish that there are no unexpected fields.
		allowedFields := stringset.New(
			"parent",        // AIP-133
			"request_id",    // AIP-155
			"validate_only", // AIP-163
			"view",          // AIP-157
			fmt.Sprintf("%s_id", strings.ToLower(strcase.SnakeCase(resourceMsgName))),
		)

		for i := 0; i < m.Fields().Len(); i++ {
			field := m.Fields().Get(i)
			// Skip the check with the field that is the body.
			if t := field.Message(); t != nil && string(t.Name()) == resourceMsgName {
				continue
			}
			// Check the remaining fields.
			if !allowedFields.Contains(string(field.Name())) && !strings.HasSuffix(string(field.Name()), "_view") {
				problems = append(problems, lint.Problem{
					Message: fmt.Sprintf(
						"Create RPCs must only contain fields explicitly described in AIPs, not %q.",
						field.Name(),
					),
					Descriptor: field,
				})
			}
		}

		return
	},
}
