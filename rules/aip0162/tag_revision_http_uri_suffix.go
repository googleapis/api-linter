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
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Tag Revision methods should have a proper HTTP pattern.
var tagRevisionHTTPURISuffix = &lint.MethodRule{
	Name:   lint.NewRuleName(162, "tag-revision-http-uri-suffix"),
	OnlyIf: utils.IsTagRevisionMethod,
	LintMethod: func(m protoreflect.MethodDescriptor) []lint.Problem {
		for _, httpRule := range utils.GetHTTPRules(m) {
			if !tagRevisionURINameRegexp.MatchString(httpRule.URI) {
				return []lint.Problem{{
					Message:    `Tag Revision URI should end with ":tagRevision".`,
					Descriptor: m,
					Location:   locations.MethodHTTPRule(m),
				}}
			}
		}

		return nil
	},
}
