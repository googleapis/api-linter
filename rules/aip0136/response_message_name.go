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

const responseMessageNameErrorMessage = "" +
	"Custom methods should return a message matching the RPC name, with a `Response` suffix, " +
	"or the resource being operated on, not %q."

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

		// Short-circuit: Output type has `Response` suffix
		suffixFindings := utils.LintMethodHasMatchingResponseName(m)
		if len(suffixFindings) == 0 {
			return nil
		}

		response := utils.GetResponseType(m)
		requestResourceType := utils.GetResourceReference(m.GetInputType().FindFieldByName("name")).GetType()
		responseResourceType := utils.GetResource(response).GetType()

		// Short-circuit: Output type is the resource being operated on
		if utils.IsResource(response) && responseResourceType == requestResourceType {
			return nil
		}

		loc := locations.MethodResponseType(m)
		if utils.IsOperation(m.GetOutputType()) {
			loc = locations.MethodOperationInfo(m)
		}

		return []lint.Problem{{
			Message:    fmt.Sprintf(responseMessageNameErrorMessage, response.GetName()),
			Descriptor: m,
			Location:   loc,
		}}

	},
}
