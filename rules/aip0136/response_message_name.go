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

	"bitbucket.org/creachadair/stringset"
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

// Custom methods should return a response message matching the RPC name,
// with a Response suffix.
var responseMessageName = &lint.MethodRule{
	Name:   lint.NewRuleName(135, "response-message-name"),
	OnlyIf: isCustomMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		matchingResponseName := m.GetName() + "Response"
		suggestion := matchingResponseName

		got := m.GetOutputType().GetName()
		want := stringset.New(matchingResponseName)

		// Get the resource name from the first request message field, if it
		// is a resource reference
		if len(m.GetInputType().GetFields()) > 0 {
			rr := utils.GetResourceReference(m.GetInputType().GetFields()[0])
			if _, resourceMsgNameSingular, ok := utils.SplitResourceTypeName(rr.GetType()); ok {
				want.Add(resourceMsgNameSingular)
				suggestion += " or " + resourceMsgNameSingular
			}
		}

		if !want.Contains(got) {
			msg := "Custom methods should return a response message matching the RPC name, with a Response suffix, or the resource, not %q."

			problem := lint.Problem{
				Message:    fmt.Sprintf(msg, got),
				Suggestion: suggestion,
				Descriptor: m,
				Location:   locations.MethodResponseType(m),
			}

			return []lint.Problem{problem}
		}

		return nil
	},
}
