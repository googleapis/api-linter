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

package aip0233

import (
	"fmt"

	"github.com/gertd/go-pluralize"
	"github.com/googleapis/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
)

// Batch Create method should have a properly named Request message.
var requestMessageName = &lint.MethodRule{
	Name:   lint.NewRuleName(233, "request-message-name"),
	OnlyIf: isBatchCreateMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		pluralInputResourceName := pluralize.NewClient().Plural(m.GetName()[11:])

		// Rule check: Establish that for methods such as `BatchCreateFoos`, the request
		// message is named `BatchCreateFoosRequest`.
		if got, want := m.GetInputType().GetName(), fmt.Sprintf("BatchCreate%sRequest", pluralInputResourceName); got != want {
			return []lint.Problem{{
				Message: fmt.Sprintf(
					"Batch Create RPCs should have a properly named request message %q, but not %q",
					want, got,
				),
				Descriptor: m,
			}}
		}

		return nil
	},
}
