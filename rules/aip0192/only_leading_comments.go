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

	"github.com/googleapis/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
)

// onlyLeadingComments ensures that a descriptor has only leading external
// comments. (Internal trailing or detached comments are permitted.)
var onlyLeadingComments = &lint.DescriptorRule{
	RuleName: lint.NewRuleName(192, "only-leading-comments"),
	LintDescriptor: func(d desc.Descriptor) []lint.Problem {
		problems := []lint.Problem{}
		c := d.GetSourceInfo()
		if len(separateInternalComments(c.GetTrailingComments()).External) > 0 {
			problems = append(problems, lint.Problem{
				Message: fmt.Sprintf(
					"%q should have only leading (not trailing) public comments.",
					d.GetName(),
				),
				Descriptor: d,
			})
		}
		if len(separateInternalComments(c.GetLeadingDetachedComments()...).External) > 0 {
			problems = append(problems, lint.Problem{
				Message: fmt.Sprintf(
					"%q has comments with empty lines between the comment and %q.",
					d.GetName(), d.GetName(),
				),
				Descriptor: d,
			})
		}
		return problems
	},
}
