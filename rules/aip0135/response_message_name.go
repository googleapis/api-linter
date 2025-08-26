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

package aip0135

import (
	"fmt"
	"strings"

	"bitbucket.org/creachadair/stringset"
	"github.com/googleapis/api-linter/v2/lint"
	"github.com/googleapis/api-linter/v2/locations"
	"github.com/googleapis/api-linter/v2/rules/internal/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Delete messages should use google.protobuf.Empty,
// google.longrunning.Operation, or the resource itself as the response
// message.
var responseMessageName = &lint.MethodRule{
	Name:   lint.NewRuleName(135, "response-message-name"),
	OnlyIf: utils.IsDeleteMethod,
	LintMethod: func(m protoreflect.MethodDescriptor) []lint.Problem {
		resource := strings.Replace(string(m.Name()), "Delete", "", 1)

		// Rule check: Establish that for methods such as `DeleteFoo`, the response
		// message is `google.protobuf.Empty` or `Foo`.
		got := m.Output().Name()
		if stringset.New("Empty", "Operation").Contains(string(got)) {
			got = protoreflect.Name(string(m.Output().FullName()))
		}
		want := stringset.New(resource, "google.protobuf.Empty")

		// If the return type is an Operation, use the annotated response type.
		lro := false
		var gotLRO string
		if utils.IsOperation(m.Output()) {
			if info := utils.GetOperationInfo(m); info != nil {
				gotLRO = info.GetResponseType()
			}
			lro = true
		}

		// If we did not get a permitted value, return a problem.
		//
		// Note: If `got` is empty string, this is an unannotated LRO.
		// The AIP-151 rule will whine about that, and this rule should not as it
		// would be confusing.
		if lro {
			if !want.Contains(gotLRO) && gotLRO != "" {
				// LRO case
				return []lint.Problem{{
					Message:    fmt.Sprintf("Delete RPCs should have response message type of Empty or the resource, not %q.", gotLRO),
					Descriptor: m,
					Location:   locations.MethodOperationInfo(m),
				}}
			}
		} else {
			if !want.Contains(string(got)) && got != "" {
				// Non-LRO case
				return []lint.Problem{{
					Message:    fmt.Sprintf("Delete RPCs should have response message type of Empty or the resource, not %q.", got),
					Suggestion: "google.protobuf.Empty",
					Descriptor: m,
					Location:   locations.MethodResponseType(m),
				}}
			}
		}

		return nil
	},
}
