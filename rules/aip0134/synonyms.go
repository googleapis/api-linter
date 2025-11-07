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

package aip0134

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/googleapis/api-linter/v2/lint"
	"github.com/googleapis/api-linter/v2/locations"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Update methods should use the word "update", not synonyms.
var synonyms = &lint.MethodRule{
	Name: lint.NewRuleName(134, "synonyms"),
	OnlyIf: func(m protoreflect.MethodDescriptor) bool {
		return m.Name() != "SetIamPolicy"
	},
	LintMethod: func(m protoreflect.MethodDescriptor) []lint.Problem {
		name := string(m.Name())
		for _, syn := range []string{"Patch", "Put", "Set"} {
			if strings.HasPrefix(name, syn) {
				synLen := len(syn)
				nameLen := len(name)
				// Check for word boundary: either exact match or next char is uppercase
				if nameLen == synLen || (nameLen > synLen && unicode.IsUpper(rune(name[synLen]))) {
					return []lint.Problem{{
						Message: fmt.Sprintf(
							`%q can be a synonym for "Update". Should this be a Update method?`,
							syn,
						),
						Descriptor: m,
						Location:   locations.DescriptorName(m),
						Suggestion: strings.Replace(name, syn, "Update", 1),
					}}
				}
			}
		}
		return nil
	},
}
