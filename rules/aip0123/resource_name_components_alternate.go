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

package aip0123

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var identifierRegexp = regexp.MustCompile("^{[a-z][_a-z0-9]*[a-z0-9]}$")

var resourceNameComponentsAlternate = &lint.MessageRule{
	Name:   lint.NewRuleName(123, "resource-name-components-alternate"),
	OnlyIf: utils.IsResource,
	LintMessage: func(m protoreflect.MessageDescriptor) []lint.Problem {
		var problems []lint.Problem
		resource := utils.GetResource(m)
		for _, p := range resource.GetPattern() {
			components := strings.Split(p, "/")
			for i, c := range components {
				identifierExpected := i%2 == 1
				if identifierExpected != isIdentifier(c) {
					problems = append(problems, lint.Problem{
						Message:    fmt.Sprintf("Resource pattern %q must alternate between collection and identifier. %q is not an identifier", p, c),
						Descriptor: m,
						Location:   locations.MessageResource(m),
					})
					break
				}
			}
		}
		return problems
	},
}

func isIdentifier(s string) bool {
	return identifierRegexp.MatchString(s)
}
