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

	"github.com/googleapis/api-linter/v2/lint"
	"github.com/googleapis/api-linter/v2/locations"
	"github.com/googleapis/api-linter/v2/rules/internal/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const responseMessageNameErrorMessage = "" +
	"Custom methods should return a message matching the RPC name, with a `Response` suffix, " +
	"or the resource being operated on, not %q."

// Custom methods should return a response message matching the RPC name,
// with a Response suffix, or the resource being operated on.
var responseMessageName = &lint.MethodRule{
	Name:   lint.NewRuleName(136, "response-message-name"),
	OnlyIf: utils.IsCustomMethod,
	LintMethod: func(m protoreflect.MethodDescriptor) []lint.Problem {
		// A response is considered valid if
		// - The response name matches the RPC name with a `Response` suffix
		// - The response is the resource being operated on
		// To identify the resource being operated on, we inspect the Resource
		// Reference of the input type's `name` field. This guidance is documented
		// in https://google.aip.dev/136#resource-based-custom-methods

		// SetIamPolicy is a special case since Policy isn't an annotated resource.
		if m.Name() == "SetIamPolicy" && m.Output().FullName() == "google.iam.v1.Policy" {
			return nil
		}

		// Short-circuit: Output type has `Response` suffix
		suffixFindings := utils.LintMethodHasMatchingResponseName(m)
		if len(suffixFindings) == 0 {
			return nil
		}

		response := utils.GetResponseType(m)
		if response == nil {
			// If the return type is not resolveable (bad) or if an LRO and
			// missing the operation_info annotation (covered by AIP-151 rules),
			// just exit.
			return nil
		}

		res := utils.GetResource(response)
		responseResourceType := res.GetType()
		requestResourceType := utils.GetResourceReference(m.Input().Fields().ByName("name")).GetType()

		// Check to see if the custom method uses the resource type name as the target
		// field name and use that instead if `name` is not present as well.
		// AIP-144 methods recommend this naming style.
		resourceFieldType := utils.GetResourceReference(m.Input().Fields().ByName(protoreflect.Name(utils.GetResourceSingular(res)))).GetType()
		if requestResourceType == "" && resourceFieldType != "" {
			requestResourceType = resourceFieldType
		}

		// Short-circuit: Output type is the resource being operated on
		if utils.IsResource(response) && responseResourceType == requestResourceType {
			return nil
		}

		loc := locations.MethodResponseType(m)
		if utils.IsOperation(m.Output()) {
			loc = locations.MethodOperationInfo(m)
		}

		return []lint.Problem{{
			Message:    fmt.Sprintf(responseMessageNameErrorMessage, response.Name()),
			Descriptor: m,
			Location:   loc,
		}}

	},
}
