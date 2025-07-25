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
	"github.com/googleapis/api-linter/locations"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// hasComments complains if there is no comment above something.
var hasComments = &lint.DescriptorRule{
	Name: lint.NewRuleName(192, "has-comments"),
	LintDescriptor: func(d protoreflect.Descriptor) (problems []lint.Problem) {
		comment := utils.SeparateInternalComments(d.ParentFile().SourceLocations().ByDescriptor(d).LeadingComments)
		if len(comment.External) == 0 {
			problems = append(problems, lint.Problem{
				Message:    fmt.Sprintf("Missing comment over %q.", d.Name()),
				Descriptor: d,
				Location:   locations.DescriptorName(d),
			})
		}
		return
	},
}
