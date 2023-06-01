// Copyright 2023 Google LLC
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

package aip0122

import (
	"strings"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

type SegmentType int64

const (
	Undefined            SegmentType = 0
	CollectionIdentifier SegmentType = 1
	ResourceId           SegmentType = 2
)

var resourceNamesAlternatingComponents = &lint.MessageRule{
	Name: lint.NewRuleName(122, "resource-names-alternating-components"),
	OnlyIf: func(m *desc.MessageDescriptor) bool {
		return utils.GetResource(m) != nil
	},
	LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
		var problems []lint.Problem
		resource := utils.GetResource(m)
		for _, p := range resource.GetPattern() {
			segs := strings.Split(p, "/")
			previousSegmentType := Undefined
			for _, seg := range segs {
				c := []rune(seg)[0]
				currentSegmentType := CollectionIdentifier
				if c == '{' {
					currentSegmentType = ResourceId
				}
				if previousSegmentType == currentSegmentType {
					return append(problems, lint.Problem{
						Message:    "Resource name components should usually alternate between collection identifiers and resource IDs.",
						Descriptor: m,
						Location:   locations.MessageResource(m),
					})
				}
				previousSegmentType = currentSegmentType
			}
		}

		return problems
	},
}
