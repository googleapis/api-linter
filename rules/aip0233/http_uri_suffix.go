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

package aip0233

import (
	"github.com/googleapis/api-linter/v2/lint"
	"github.com/googleapis/api-linter/v2/locations"
	"github.com/googleapis/api-linter/v2/rules/internal/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Batch Create methods should have a proper HTTP pattern.
var httpURISuffix = &lint.MethodRule{
	Name:   lint.NewRuleName(233, "http-uri-suffix"),
	OnlyIf: isBatchCreateMethod,
	LintMethod: func(m protoreflect.MethodDescriptor) []lint.Problem {
		for _, httpRule := range utils.GetHTTPRules(m) {
			if !batchCreateURINameRegexp.MatchString(httpRule.URI) {
				return []lint.Problem{{
					Message:    `Batch Create methods URI should end with ":batchCreate".`,
					Descriptor: m,
					Location:   locations.MethodHTTPRule(m),
				}}
			}
		}

		return nil
	},
}
