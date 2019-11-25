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

package aip0136

import (
	"strings"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/jhump/protoreflect/desc"
	"github.com/stoewer/go-strcase"
)

var verbNoun = &descrule.MethodRule{
	RuleName: lint.NewRuleName(136, "verb-noun"),
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		// We can not detect this precisely without a full dictionary (probably
		// not worth it), but we can catch some common mistakes.

		// Are there at least two words?
		if wordCount := len(strings.Split(strcase.SnakeCase(m.GetName()), "_")); wordCount == 1 {
			return []lint.Problem{{
				Message:    "Custom methods should be named using a verb followed by a noun.",
				Descriptor: m,
				Location:   locations.DescriptorName(m),
			}}
		}

		return nil
	},
}
