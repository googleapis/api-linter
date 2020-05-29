// Copyright 2020 Google LLC
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

package aip0140

import (
	"strings"

	"bitbucket.org/creachadair/stringset"
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/jhump/protoreflect/desc"
)

var uri = &lint.FieldRule{
	Name: lint.NewRuleName(140, "uri"),
	LintField: func(f *desc.FieldDescriptor) []lint.Problem {
		nameSegments := stringset.New(strings.Split(f.GetName(), "_")...)
		if nameSegments.Contains("url") {
			return []lint.Problem{{
				Message:    "Use `uri` instead of `url` in field names.",
				Descriptor: f,
				Location:   locations.DescriptorName(f),
				Suggestion: strings.ReplaceAll(f.GetName(), "url", "uri"),
			}}
		}
		return nil
	},
}
