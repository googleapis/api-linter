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

package aip0192

import (
	"regexp"
	"strings"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/jhump/protoreflect/desc"
)

// absoluteLinks ensures that a descriptor has only absolute links.
var absoluteLinks = &lint.DescriptorRule{
	Name: lint.NewRuleName(192, "absolute-links"),
	LintDescriptor: func(d desc.Descriptor) []lint.Problem {
		comment := strings.Join(
			separateInternalComments(d.GetSourceInfo().GetLeadingComments()).External,
			"\n",
		)

		// Absolute links can appear in either free text or Markdown links,
		// so we need to spot both, but we also do not want a ton of false
		// positives.
		//
		// It seems infeasible to check free text, but we can check Markdown links.
		for _, match := range mdLink.FindAllStringSubmatch(comment, -1) {
			uri := match[1]
			if !strings.Contains(uri, "://") {
				return []lint.Problem{{
					Message:    "Links in comments should be absolute and begin with https://",
					Descriptor: d,
					Location:   locations.DescriptorName(d),
				}}
			}
		}

		return nil
	},
}

var mdLink = regexp.MustCompile(`\[.*?\]\(([\w\d:/.%_-]+)\)`)
