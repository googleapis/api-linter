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

package aip0131

import (
	"fmt"
	"strings"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/jhump/protoreflect/desc"
)

// Get methods should not generally use synonyms for "get".
var synonyms = &lint.MethodRule{
	Name: lint.NewRuleName(131, "synonyms"),
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		name := m.GetName()
		for _, syn := range []string{"Acquire", "Fetch", "Lookup", "Read", "Retrieve"} {
			if strings.HasPrefix(name, syn) {
				return []lint.Problem{{
					Message: fmt.Sprintf(
						`%q can be a synonym for "Get". Should this be a Get method?`,
						syn,
					),
					Suggestion: strings.Replace(name, syn, "Get", 1),
					Descriptor: m,
					Location:   locations.DescriptorName(m),
				}}
			}
		}
		return nil
	},
}
