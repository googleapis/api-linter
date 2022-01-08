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

	"bitbucket.org/creachadair/stringset"
	"github.com/commure/api-linter/lint"
	"github.com/commure/api-linter/locations"
	"github.com/commure/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
	"github.com/stoewer/go-strcase"
)

var httpBody = &lint.MethodRule{
	Name:   lint.NewRuleName(136, "http-body"),
	OnlyIf: isCustomMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		for _, httpRule := range utils.GetHTTPRules(m) {
			noBody := stringset.New("GET", "DELETE")
			if !noBody.Contains(httpRule.Method) {
				// Determine the name of the resource.
				// This entails some guessing; we assume that the verb is a single
				// word and that the resource is everything else.
				resource := strings.Join(strings.Split(strcase.SnakeCase(m.GetName()), "_")[1:], "_")
				if !stringset.New(resource, "*").Contains(httpRule.Body) {
					// FIXME: We intentionally only return one Problem here.
					// When we can attach problems to the particular annotation, update this
					// to return multiples.
					return []lint.Problem{{
						Message:    "Custom POST methods should set `body: \"*\"`.",
						Descriptor: m,
						Location:   locations.MethodHTTPRule(m),
					}}
				}
			} else if noBody.Contains(httpRule.Method) && httpRule.Body != "" {
				// FIXME: We intentionally only return one Problem here.
				// When we can attach problems to the particular annotation, update this
				// to return multiples.
				return []lint.Problem{{
					Message:    "Custom GET (or DELETE) methods should not set a body clause.",
					Descriptor: m,
					Location:   locations.MethodHTTPRule(m),
				}}
			}
		}
		return nil
	},
}
