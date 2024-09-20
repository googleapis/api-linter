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
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

var tagRevisionResponseMessageName = &lint.MethodRule{
	Name:   lint.NewRuleName(162, "tag-revision-response-message-name"),
	OnlyIf: utils.IsTagRevisionMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		// Rule check: Establish that for methods such as `TagBookRevision`, the response
		// message is `Book`.
		want, ok := utils.ExtractRevisionResource(m)
		if !ok {
			return nil
		}
		response := utils.GetResponseType(m)
		if response == nil {
			return nil
		}
		got := response.GetName()
		loc := locations.MethodResponseType(m)
		suggestion := want

		if utils.GetOperationInfo(m) != nil {
			loc = locations.MethodOperationInfo(m)
			suggestion = "" // We cannot offer a precise enough location to make a suggestion.
		}

		// Return a problem if we did not get the expected return name.
		if got != want {
			return []lint.Problem{{
				Message: fmt.Sprintf(
					"Tag Revision RPCs should have response message type %q, not %q.",
					want,
					got,
				),
				Suggestion: suggestion,
				Descriptor: m,
				Location:   loc,
			}}
		}
		return nil
	},
}
