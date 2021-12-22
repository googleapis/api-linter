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

	"github.com/googleapis/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
)

var requestNameRequired = &lint.MessageRule{
	Name:   lint.NewRuleName(135, "request-name-required"),
	OnlyIf: isDeleteRequestMessage,
	LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
		if m.FindFieldByName("resource_name") == nil {
			return []lint.Problem{{
				Message:    fmt.Sprintf("Method %q has no `resource_name` field", m.GetName()),
				Descriptor: m,
			}}
		}
		return nil
	},
}
