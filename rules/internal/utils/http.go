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
	"regexp"

	"github.com/jhump/protoreflect/desc"
	apb "google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
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
	if x := proto.GetExtension(opts, apb.E_Http); x != nil {
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
	oneof := map[string]string{
		"GET":    rule.GetGet(),
		"POST":   rule.GetPost(),
		"PUT":    rule.GetPut(),
		"PATCH":  rule.GetPatch(),
		"DELETE": rule.GetDelete(),
	}
	if custom := rule.GetCustom(); custom != nil {
		oneof[custom.GetKind()] = custom.GetPath()
	}
	for method, uri := range oneof {
		if uri != "" {
			return &HTTPRule{
				Method:       method,
				URI:          uri,
				Body:         rule.GetBody(),
				ResponseBody: rule.GetResponseBody(),
			}
		}
	}
	return nil
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

// GetVariables returns the variable segments in a URI as a map.
func (h *HTTPRule) GetVariables() map[string]string {
	vars := map[string]string{}
	for _, match := range plainVar.FindAllStringSubmatch(h.URI, -1) {
		vars[match[1]] = "*"
	}
	for _, match := range varSegment.FindAllStringSubmatch(h.URI, -1) {
		vars[match[1]] = match[2]
	}
	return vars
}

// GetPlainURI returns the URI with variable segment information removed.
func (h *HTTPRule) GetPlainURI() string {
	return plainVar.ReplaceAllString(varSegment.ReplaceAllString(h.URI, "$2"), "*")
}

var plainVar = regexp.MustCompile(`\{([^}=]+)}`)
var varSegment = regexp.MustCompile(`\{([^}=]+)=([^}]+)\}`)
