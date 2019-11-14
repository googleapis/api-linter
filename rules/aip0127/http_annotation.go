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

package aip0127

import (
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

var hasAnnotation = &lint.MethodRule{
	Name: lint.NewRuleName(127, "http-annotation"),
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		hasHTTPRule := len(utils.GetHTTPRules(m)) > 0
		if hasHTTPRule && m.IsClientStreaming() && m.IsServerStreaming() {
			return []lint.Problem{{
				Message:    "Bi-directional streaming RPCs should omit `google.api.http`.",
				Descriptor: m,
			}}
		}
		if !hasHTTPRule && !(m.IsClientStreaming() && m.IsServerStreaming()) {
			return []lint.Problem{{
				Message:    "RPCs must include HTTP definitions using the `google.api.http` annotation.",
				Descriptor: m,
			}}
		}
		return nil
	},
}
