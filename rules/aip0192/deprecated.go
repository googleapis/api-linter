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

package aip0192

import (
	"strings"

	"github.com/googleapis/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
)

var deprecatedMethodComment = &lint.MethodRule{
	Name:   lint.NewRuleName(192, "deprecated-method-comment"),
	OnlyIf: isDeprecatedMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		comment := separateInternalComments(m.GetSourceInfo().GetLeadingComments())
		if len(comment.External) > 0 && strings.HasPrefix(comment.External[0], "Deprecated:") {
			return nil
		}
		return []lint.Problem{{
			Message:    "Deprecated methods must have \"Deprecated: <reason>\" as the first line in the method comment.",
			Descriptor: m,
		}}
	},
}

var deprecatedServiceComment = &lint.ServiceRule{
	Name:   lint.NewRuleName(192, "deprecated-service-comment"),
	OnlyIf: isDeprecatedService,
	LintService: func(s *desc.ServiceDescriptor) []lint.Problem {
		comment := separateInternalComments(s.GetSourceInfo().GetLeadingComments())
		if len(comment.External) > 0 && strings.HasPrefix(comment.External[0], "Deprecated:") {
			return nil
		}

		return []lint.Problem{{
			Message:    "Deprecated services must have \"Deprecated: <reason>\" as the first line in the service comment.",
			Descriptor: s,
		}}
	},
}
