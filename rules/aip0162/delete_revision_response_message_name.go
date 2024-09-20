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

// Delete Revision methods should return the resource itself.
var deleteRevisionResponseMessageName = &lint.MethodRule{
	Name:   lint.NewRuleName(162, "delete-revision-response-message-name"),
	OnlyIf: utils.IsDeleteRevisionMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
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

		if got != want {
			return []lint.Problem{{
				Message:    fmt.Sprintf("Delete Revision methods should return the resource itself (`%s`), not `%s`.", want, got),
				Suggestion: suggestion,
				Descriptor: m,
				Location:   loc,
			}}
		}
		return nil
	},
}
