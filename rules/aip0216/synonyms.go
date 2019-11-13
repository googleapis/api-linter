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

package aip0216

import (
	"strings"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/jhump/protoreflect/desc"
)

var synonyms = &lint.EnumRule{
	Name: lint.NewRuleName(216, "synonyms"),
	LintEnum: func(e *desc.EnumDescriptor) []lint.Problem {
		if strings.HasSuffix(e.GetName(), "Status") {
			return []lint.Problem{{
				Message:    `Prefer "State" over "Status" for lifecycle state enums.`,
				Suggestion: strings.Replace(e.GetName(), "Status", "State", 1),
				Descriptor: e,
				Location:   locations.DescriptorName(e),
			}}
		}
		return nil
	},
}
