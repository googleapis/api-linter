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
	"github.com/googleapis/api-linter/v2/lint"
	"github.com/googleapis/api-linter/v2/locations"
	"github.com/googleapis/api-linter/v2/rules/internal/utils"
	"github.com/stoewer/go-strcase"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var allowedBodyTypes = stringset.New("google.api.HttpBody")

var httpBody = &lint.MethodRule{
	Name:   lint.NewRuleName(136, "http-body"),
	OnlyIf: utils.IsCustomMethod,
	LintMethod: func(m protoreflect.MethodDescriptor) []lint.Problem {
		for _, httpRule := range utils.GetHTTPRules(m) {
			noBody := stringset.New("GET", "DELETE")
			if !noBody.Contains(httpRule.Method) {
				// Determine the name of the resource.
				// This entails some guessing; we assume that the verb is a single
				// word and that the resource is everything else.
				resource := strings.Join(strings.Split(strcase.SnakeCase(string(m.Name())), "_")[1:], "_")
				if !stringset.New(resource, "*").Contains(httpRule.Body) {
					// Some field types make sense to use as a body instead of the entire request message,
					// e.g. google.api.HttpBody.
					b := m.Input().Fields().ByName(protoreflect.Name(httpRule.Body))
					isMessageType := b != nil && b.Message() != nil
					if isMessageType && allowedBodyTypes.Contains(string(b.Message().FullName())) {
						continue
					}

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
