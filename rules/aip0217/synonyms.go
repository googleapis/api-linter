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

package aip0217

import (
	"bitbucket.org/creachadair/stringset"
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var synSet = stringset.New(
	"unreachable_locations",
)

var synonyms = &lint.FieldRule{
	Name: lint.NewRuleName(217, "synonyms"),
	OnlyIf: func(f protoreflect.FieldDescriptor) bool {
		return f.GetOwner().FindFieldByName("next_page_token") != nil
	},
	LintField: func(f protoreflect.FieldDescriptor) []lint.Problem {
		if synSet.Contains(f.Name()) {
			return []lint.Problem{{
				Message:    "Use `unreachable` to express unreachable resources when listing.",
				Suggestion: "unreachable",
				Descriptor: f,
				Location:   locations.DescriptorName(f),
			}}
		}
		return nil
	},
}
