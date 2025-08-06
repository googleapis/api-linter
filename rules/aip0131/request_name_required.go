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

package aip0131

import (
	"fmt"

	"github.com/googleapis/api-linter/v2/lint"
	"github.com/googleapis/api-linter/v2/rules/internal/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// The Get standard method should have some required fields.
var requestNameRequired = &lint.MessageRule{
	Name:   lint.NewRuleName(131, "request-name-required"),
	OnlyIf: utils.IsGetRequestMessage,
	LintMessage: func(m protoreflect.MessageDescriptor) []lint.Problem {
		if m.Fields().ByName("name") == nil {
			return []lint.Problem{{
				Message:    fmt.Sprintf("Method %q has no `name` field", m.Name()),
				Descriptor: m,
			}}
		}

		return nil
	},
}
