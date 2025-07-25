// Copyright 2020 Google LLC
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

package aip0154

import (
	"fmt"
	"strings"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var declarativeFriendlyRequired = &lint.MessageRule{
	Name: lint.NewRuleName(154, "declarative-friendly-required"),
	OnlyIf: func(m protoreflect.MessageDescriptor) bool {
		// Sanity check: If the resource is not declarative-friendly, none of
		// this logic applies.
		if resource := utils.DeclarativeFriendlyResource(m); resource != nil {
			// This should apply if the resource in question is declarative-friendly,
			// but our IsDeclarativeFriendly method will return true for both
			// resources and request messages, and they need to be handled subtly
			// differently.
			if m == resource {
				return true
			}

			// If this is a request message, then make several more checks based on
			// what the method looks like.
			if name := string(m.Name()); strings.HasSuffix(name, "Request") {
				methodName := strings.TrimSuffix(name, "Request")

				// If this is a GET request, then this message is exempt.
				if method := utils.FindMethod(m.Parent().(protoreflect.FileDescriptor), methodName); method != nil {
					for _, rule := range utils.GetHTTPRules(method) {
						if rule.Method == "GET" {
							return false
						}
					}
				}

				// If the message contains the resource, then this message is exempt.
				for i := 0; i < m.Fields().Len(); i++ {
					field := m.Fields().Get(i)
					if field.Message() == resource {
						return false
					}
				}

				// Okay, this message should include an etag.
				return true
			}
		}

		return false
	},
	LintMessage: func(m protoreflect.MessageDescriptor) []lint.Problem {
		for i := 0; i < m.Fields().Len(); i++ {
			field := m.Fields().Get(i)
			if field.Name() == "etag" {
				return nil
			}
		}

		whoami := "resources"
		if strings.HasSuffix(string(m.Name()), "Request") {
			whoami = "mutation requests without the resource"
		}
		return []lint.Problem{{
			Message:    fmt.Sprintf("Declarative-friendly %s should include `string etag`.", whoami),
			Descriptor: m,
		}}
	},
}
