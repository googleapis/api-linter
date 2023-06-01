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
	"regexp"
	"strings"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

var noMarkdownHeadings = &lint.DescriptorRule{
	Name: lint.NewRuleName(192, "no-markdown-headings"),
	LintDescriptor: func(d desc.Descriptor) []lint.Problem {
		for _, cmt := range utils.SeparateInternalComments(d.GetSourceInfo().GetLeadingComments()).External {
			for _, line := range strings.Split(cmt, "\n") {
				if heading.FindString(line) != "" {
					return []lint.Problem{{
						Message:    "Comments should not include Markdown headings.",
						Descriptor: d,
					}}
				}
			}
		}
		return nil
	},
}

var heading = regexp.MustCompile(`^\s*[#]+ `)
