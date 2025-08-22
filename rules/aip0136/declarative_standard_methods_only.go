// Copyright 2020 Google LLC
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

package aip0136

import (
	"strings"

	"bitbucket.org/creachadair/stringset"
	"github.com/googleapis/api-linter/v2/lint"
	"github.com/googleapis/api-linter/v2/rules/internal/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var standardMethodsOnly = &lint.MethodRule{
	Name:   lint.NewRuleName(136, "declarative-standard-methods-only"),
	OnlyIf: utils.IsDeclarativeFriendlyMethod,
	LintMethod: func(m protoreflect.MethodDescriptor) []lint.Problem {
		// Standard methods are fine.
		standard := stringset.New("Get", "List", "Create", "Update", "Delete", "Undelete", "Batch")
		for s := range standard {
			if strings.HasPrefix(string(m.Name()), s) {
				return nil
			}
		}

		// This is likely to have a non-trivial number of exceptions, and a
		// traditional linter disable may not be appropriate.
		//
		// Therefore, we allow "Imperative only." in an internal comment to make
		// this not complain.
		if cmt := m.ParentFile().SourceLocations().ByDescriptor(m).LeadingComments; strings.Contains(strings.ToLower(string(cmt)), "imperative only") {
			return nil
		}

		// Okay, complain.
		return []lint.Problem{{
			Message: strings.Join([]string{
				"Declarative-friendly resources should generally avoid custom methods.\n",
				"However, if this is an imperative-only method that does *not* need ",
				`declarative tooling support, add the text "Imperative only." to the comment. `,
				"(Using an internal comment is fine.)",
			}, ""),
			Descriptor: m,
		}}
	},
}