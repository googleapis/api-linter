// Copyright 2026 Google LLC
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

package aip0190

import (
	"fmt"

	"github.com/googleapis/api-linter/v2/lint"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var messageCase = &lint.MessageRule{
	Name: lint.NewRuleName(190, "message-case"),
	LintMessage: func(m protoreflect.MessageDescriptor) []lint.Problem {
		name := string(m.Name())
		if !isValidCamelCase(name) {
			return []lint.Problem{{
				Message:    fmt.Sprintf("Message name %q must use UpperCamelCase.", name),
				Descriptor: m,
			}}
		}
		return nil
	},
}
