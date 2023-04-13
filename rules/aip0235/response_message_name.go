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

package aip0235

import (
	"fmt"

	"bitbucket.org/creachadair/stringset"
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

// Batch Delete method should have a properly named response message.
var responseMessageName = &lint.MethodRule{
	Name:   lint.NewRuleName(235, "response-message-name"),
	OnlyIf: isBatchDeleteMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		got := m.GetOutputType().GetFullyQualifiedName()
		if utils.IsOperation(m.GetOutputType()) {
			got = utils.GetOperationInfo(m).GetResponseType()
		} else if got != "google.protobuf.Empty" {
			got = m.GetOutputType().GetName()
		}

		wantSoftDelete := m.GetName() + "Response"
		want := stringset.New(
			"google.protobuf.Empty",
			wantSoftDelete,
		)

		// Rule check: Establish that for methods such as `BatchDeleteFoos`, the
		// response message should be named `BatchDeleteFoosResponse` or
		// `google.protobuf.Empty`.
		//
		// Note: If `got` is empty string, this is an unannotated LRO.
		// The AIP-151 rule will whine about that, and this rule should not as it
		// would be confusing.
		if !want.Contains(got) && got != "" {
			return []lint.Problem{{
				Message: fmt.Sprintf(
					"Batch Delete RPCs should have response message type of `google.protobuf.Empty` or `%s` (for soft-deletes only), not %q.",
					wantSoftDelete, got,
				),
				Suggestion: "google.protobuf.Empty",
				Descriptor: m,
				Location:   locations.MethodResponseType(m),
			}}
		}

		return nil
	},
}
