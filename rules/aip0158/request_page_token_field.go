// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aip0158

import (
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

var requestPaginationPageToken = &lint.MessageRule{
	Name:   lint.NewRuleName(158, "request-page-token-field"),
	OnlyIf: isPaginatedRequestMessage,
	LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
		f, problems := utils.LintFieldPresent(m, "page_token")
		if len(problems) > 0 {
			return problems
		}
		// Checks that page_size is of type string and is not a oneof. These are
		// noops if page_size is not a oneof and is a string.
		problems = append(problems, utils.LintSingularStringField(f)...)
		problems = append(problems, utils.LintNotOneof(f)...)

		return problems
	},
}
