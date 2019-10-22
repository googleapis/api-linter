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

package utils

import (
	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/desc"
	apb "google.golang.org/genproto/googleapis/api/annotations"
)

// GetHTTPRules returns a slice of HTTP rules for a given method descriptor.
//
// Note: This returns a slice -- it takes the google.api.http annotation,
// and then flattens the values in `additional_bindings`.
// This allows rule authors to simply range over all of the HTTP rules,
// since the common case is to want to apply the checks to all of them.
func GetHTTPRules(m *desc.MethodDescriptor) []*HTTPRule {
	rules := []*HTTPRule{}

	// Get the method options.
	opts := m.GetMethodOptions()

	// Get the "primary" rule (the direct google.api.http annotation).
	if x, err := proto.GetExtension(opts, apb.E_Http); err == nil {
		httpRule := x.(*apb.HttpRule)
		if parsedRule := parseRule(httpRule); parsedRule != nil {
			rules = append(rules, parsedRule)

			// Add any additional bindings and flatten them into `rules`.
			for _, binding := range httpRule.GetAdditionalBindings() {
				rules = append(rules, parseRule(binding))
			}
		}
	}

	// Done; return the rules.
	return rules
}

func parseRule(rule *apb.HttpRule) *HTTPRule {
	answer := &HTTPRule{}
	if uri := rule.GetGet(); uri != "" {
		answer.Method = "GET"
		answer.URI = uri
	} else if uri := rule.GetPost(); uri != "" {
		answer.Method = "POST"
		answer.URI = uri
	} else if uri := rule.GetPut(); uri != "" {
		answer.Method = "PUT"
		answer.URI = uri
	} else if uri := rule.GetPatch(); uri != "" {
		answer.Method = "PATCH"
		answer.URI = uri
	} else if uri := rule.GetDelete(); uri != "" {
		answer.Method = "DELETE"
		answer.URI = uri
	} else if custom := rule.GetCustom(); custom != nil {
		answer.Method = custom.GetKind()
		answer.URI = custom.GetPath()
	}

	// Set the body and response body, and return the answer.
	answer.Body = rule.GetBody()
	answer.ResponseBody = rule.GetResponseBody()
	return answer
}

// HTTPRule defines a parsed, easier-to-query equivalent to `apb.HttpRule`.
type HTTPRule struct {
	// The HTTP method. Guaranteed to be in all caps.
	// This is set to "CUSTOM" if the Custom property is set.
	Method string

	// The HTTP URI (the value corresponding to the selected HTTP method).
	URI string

	// The `body` value forwarded from the generated proto's HttpRule.
	Body string

	// The `response_body` value forwarded from the generated proto's HttpRule.
	ResponseBody string
}
