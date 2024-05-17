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

package aip0136

import (
	"fmt"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

// Custom methods should return a response message matching the RPC name,
// with a Response suffix, or the resource being operated on.
var responseMessageName = &lint.MethodRule{
	Name:   lint.NewRuleName(136, "response-message-name"),
	OnlyIf: isCustomMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		// A response is considered valid if
		// - The response name matches the RPC name with a `Response` suffix
		// - The response is the resource being operated on
		// To identify the resource being operated on, we inspect the Resource
		// Reference of the input type's `name` field. This guidance is documented
		// in https://google.aip.dev/136#resource-based-custom-methods

		response := utils.GetResponseType(m)
		requestType := utils.GetResourceReference(m.GetInputType().FindFieldByName("name")).GetType()
		resourceOperatedOn := utils.FindResourceMessage(requestType, m.GetFile())

		// Output type has `Response` suffix
		if len(utils.LintMethodHasMatchingResponseName(m)) == 0 {
			return nil
		}

		// Response is the resource being operated on
		if response == resourceOperatedOn {
			return nil
		}

		msg := "Custom methods should return a message matching the RPC name, with a `Response` suffix, or the resource being operated on, not %q."
		got := m.GetName()

		suggestion := got + "Response"
		if utils.IsResource(response) && resourceOperatedOn != nil {
			suggestion += " or " + resourceOperatedOn.GetName()
		}

		return []lint.Problem{{
			Message:    fmt.Sprintf(msg, got),
			Suggestion: suggestion,
			Descriptor: m,
			Location:   locations.MethodResponseType(m),
		}}

	},
}
