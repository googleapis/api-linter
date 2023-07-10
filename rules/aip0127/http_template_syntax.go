// Copyright 2022 Google LLC
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
	"fmt"
	"regexp"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

var (
	literal          = `(\w+)`
	identifier       = `(\w+)`
	verb             = fmt.Sprintf("(:%s)", literal)
	fieldPath        = fmt.Sprintf(`(%s(\.%s)*)`, identifier, identifier)
	variableSegment  = fmt.Sprintf(`(\*\*|\*|%s)`, literal)
	variableTemplate = fmt.Sprintf("(%s(/%s)*)", variableSegment, variableSegment)
	variable         = fmt.Sprintf(`(\{%s(=%s)?\})`, fieldPath, variableTemplate)
	segment          = fmt.Sprintf(`(\*\*|\*|%s|%s)`, literal, variable)
	segments         = fmt.Sprintf("(%s(/%s)*)", segment, segment)
	template         = fmt.Sprintf("^/%s(%s)?$", segments, verb)
	templateRegex    = regexp.MustCompile(template)
)

// HTTP URL pattern should follow the syntax rules described here:
// https://github.com/googleapis/googleapis/blob/16db2fb7fab4668bdfa09966513e03581d8f5e35/google/api/http.proto#L224.
var httpTemplateSyntax = &lint.MethodRule{
	Name:   lint.NewRuleName(127, "http-template-syntax"),
	OnlyIf: utils.HasHTTPRules,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		problems := []lint.Problem{}
		for _, httpRule := range utils.GetHTTPRules(m) {
			// Replace the API Versioning template if it matches exactly so as
			// to not emit false positives.
			uri := utils.VersionedSegment.ReplaceAllString(httpRule.URI, "v")
			if !templateRegex.MatchString(uri) {
				message := fmt.Sprintf("The HTTP pattern %q does not follow proper HTTP path template syntax", httpRule.URI)
				problems = append(problems, lint.Problem{
					Message:    message,
					Descriptor: m,
					Location:   locations.MethodHTTPRule(m),
				})
			}
		}
		return problems
	},
}
