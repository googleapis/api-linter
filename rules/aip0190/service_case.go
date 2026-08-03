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

var serviceCase = &lint.ServiceRule{
	Name: lint.NewRuleName(190, "service-case"),
	LintService: func(s protoreflect.ServiceDescriptor) []lint.Problem {
		name := string(s.Name())
		if !isValidCamelCase(name) {
			return []lint.Problem{{
				Message:    fmt.Sprintf("Service name %q must use UpperCamelCase.", name),
				Descriptor: s,
			}}
		}
		return nil
	},
}
