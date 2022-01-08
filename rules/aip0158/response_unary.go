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

package aip0158

import (
	"github.com/commure/api-linter/lint"
	"github.com/commure/api-linter/locations"
	"github.com/jhump/protoreflect/desc"
)

var responseUnary = &lint.MethodRule{
	Name:   lint.NewRuleName(158, "response-unary"),
	OnlyIf: isPaginatedMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		if m.IsServerStreaming() {
			return []lint.Problem{{
				Message:    "Paginated responses must be unary, not streaming.",
				Descriptor: m,
				Location:   locations.MethodResponseType(m),
			}}
		}
		return nil
	},
}
