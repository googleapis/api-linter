// Copyright 2024 Google LLC
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

package aip0217

import (
	"github.com/googleapis/api-linter/lint"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var returnPartialSuccessWithUnreachable = &lint.MethodRule{
	Name: lint.NewRuleName(217, "return-partial-success-with-unreachable"),
	OnlyIf: func(m protoreflect.MethodDescriptor) bool {
		return m.Input().Fields().ByName("return_partial_success") != nil
	},
	LintMethod: func(m protoreflect.MethodDescriptor) (problems []lint.Problem) {
		if m.Output().Fields().ByName("unreachable") == nil {
			problems = append(problems, lint.Problem{
				Message:    "`return_partial_success` must be added in conjunction with response field `repeated string unreachable`.",
				Descriptor: m.Input().Fields().ByName("return_partial_success"),
			})
		}
		return
	},
}
