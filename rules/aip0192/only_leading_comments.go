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

package aip0192

import (
	"fmt"

	"github.com/googleapis/api-linter/v2/lint"
	"github.com/googleapis/api-linter/v2/rules/internal/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// onlyLeadingComments ensures that a descriptor has only leading external
// comments. (Internal trailing or detached comments are permitted.)
var onlyLeadingComments = &lint.DescriptorRule{
	Name: lint.NewRuleName(192, "only-leading-comments"),
	LintDescriptor: func(d protoreflect.Descriptor) []lint.Problem {
		problems := []lint.Problem{}
		if len(utils.SeparateInternalComments(d.ParentFile().SourceLocations().ByDescriptor(d).TrailingComments).External) > 0 {
			problems = append(problems, lint.Problem{
				Message: fmt.Sprintf(
					"%q should have only leading (not trailing) public comments.",
					d.Name(),
				),
				Descriptor: d,
			})
		}
		if len(utils.SeparateInternalComments(d.ParentFile().SourceLocations().ByDescriptor(d).LeadingDetachedComments...).External) > 0 {
			problems = append(problems, lint.Problem{
				Message: fmt.Sprintf(
					"%q has comments with empty lines between the comment and %q.",
					d.Name(), d.Name(),
				),
				Descriptor: d,
			})
		}
		return problems
	},
}
