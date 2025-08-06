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

	"github.com/googleapis/api-linter/v2/lint"
	"github.com/googleapis/api-linter/v2/locations"
	"github.com/googleapis/api-linter/v2/rules/internal/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
	"github.com/stoewer/go-strcase"
)

// The create request message should have resource field.
var resourceField = &lint.MessageRule{
	Name:   lint.NewRuleName(133, "request-resource-field"),
	OnlyIf: utils.IsCreateRequestMessage,
	LintMessage: func(m protoreflect.MessageDescriptor) []lint.Problem {
		resourceMsgName := getResourceMsgNameFromReq(m)

		// The rule (resource field name must map to the POST body) is
		// checked by AIP-0133 ("core::0133::http-body")
		for i := 0; i < m.Fields().Len(); i++ {
			fieldDesc := m.Fields().Get(i)
			if msgDesc := fieldDesc.Message(); msgDesc != nil && string(msgDesc.Name()) == resourceMsgName {
				// Rule check: Is the field named properly?
				if want := strcase.SnakeCase(resourceMsgName); string(fieldDesc.Name()) != want {
					return []lint.Problem{{
						Message:    fmt.Sprintf("Resource field should be named %q.", want),
						Descriptor: fieldDesc,
						Suggestion: want,
						Location:   locations.DescriptorName(fieldDesc),
					}}
				}
				return nil
			}
		}

		// Rule check: Establish that a resource field must be included.
		return []lint.Problem{{
			Message:    fmt.Sprintf("Message %q has no %q type field", m.Name(), resourceMsgName),
			Descriptor: m,
		}}
	},
}