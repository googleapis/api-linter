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
	"strings"

	"github.com/commure/api-linter/lint"
	"github.com/commure/api-linter/locations"
	"github.com/jhump/protoreflect/desc"
)

// Delete Revision methods should return the resource itself.
var deleteRevisionResponseMessageName = &lint.MethodRule{
	Name:   lint.NewRuleName(162, "delete-revision-response-message-name"),
	OnlyIf: isDeleteRevisionMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		got := m.GetOutputType().GetName()
		want := strings.TrimPrefix(strings.TrimSuffix(m.GetName(), "Revision"), "Delete")
		if got != want {
			return []lint.Problem{{
				Message:    fmt.Sprintf("Delete Revision methods should return the resource itself (`%s`), not `%s`.", want, got),
				Suggestion: want,
				Descriptor: m,
				Location:   locations.MethodResponseType(m),
			}}
		}
		return nil
	},
}
