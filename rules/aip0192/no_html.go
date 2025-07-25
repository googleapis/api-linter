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
	"regexp"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var noHTML = &lint.DescriptorRule{
	Name: lint.NewRuleName(192, "no-html"),
	OnlyIf: func(d protoreflect.Descriptor) bool {
		return len(d.ParentFile().SourceLocations().ByDescriptor(d).LeadingComments) > 0
	},
	LintDescriptor: func(d protoreflect.Descriptor) []lint.Problem {
		for _, comment := range utils.SeparateInternalComments(d.ParentFile().SourceLocations().ByDescriptor(d).LeadingComments).External {
			if htmlTag.MatchString(comment) {
				return []lint.Problem{{
					Message:    "Comments must not include raw HTML.",
					Descriptor: d,
				}}
			}
		}
		return nil
	},
}

// Yes, yes, I know: https://stackoverflow.com/questions/1732348/
// Even Jon Skeet cannot parse HTML using regular expressions.
//
// That said, we really only want to pick up "basic HTML smell", and are
// disinterested in actually doing any manipulation, and we can be at least
// a little tolerant of false positives/negatives. To do this, we'll look for
// closing tags (e.g., <foo /> or </foo>) since opening tags are likely to be
// used for variable placeholders (e.g., http://<server>/<path>).
//
// TL;DR: In this case, a regex seems better than taking a dependency just for
// checking if HTML exists in comments.
var htmlTag = regexp.MustCompile(`(</ *[a-zA-Z-]+>|<[a-zA-Z-]+ */>)`)
