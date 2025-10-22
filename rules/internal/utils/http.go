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
	"strings"

	apb "google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// HasHTTPRules returns true when the given method descriptor is annotated with
// a google.api.http option.
func HasHTTPRules(m protoreflect.MethodDescriptor) bool {
	return m.Options().ProtoReflect().Has(apb.E_Http.TypeDescriptor())
}

// GetHTTPRules returns a slice of HTTP rules for a given method descriptor.
//
// Note: This returns a slice -- it takes the google.api.http annotation,
// and then flattens the values in `additional_bindings`.
// This allows rule authors to simply range over all of the HTTP rules,
// since the common case is to want to apply the checks to all of them.
func GetHTTPRules(m protoreflect.MethodDescriptor) []*HTTPRule {
	rules := []*HTTPRule{}
	opts := m.Options()
	if !opts.ProtoReflect().Has(apb.E_Http.TypeDescriptor()) {
		return rules
	}

	// Get the "primary" rule (the direct google.api.http annotation).
	ext := opts.ProtoReflect().Get(apb.E_Http.TypeDescriptor()).Message()
	if parsedRule := parseRule(ext); parsedRule != nil {
		rules = append(rules, parsedRule)
	}

	// Add any additional bindings.
	additionalBindingsDesc := ext.Descriptor().Fields().ByName("additional_bindings")
	if additionalBindingsDesc != nil {
		bindings := ext.Get(additionalBindingsDesc).List()
		for i := 0; i < bindings.Len(); i++ {
			if parsedRule := parseRule(bindings.Get(i).Message()); parsedRule != nil {
				rules = append(rules, parsedRule)
			}
		}
	}

	// Done; return the rules.
	return rules
}

func parseRule(rule protoreflect.Message) *HTTPRule {
	var method, uri, body, responseBody string
	rule.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		switch fd.Name() {
		case "body":
			body = v.String()
		case "response_body":
			responseBody = v.String()
		case "get", "put", "post", "delete", "patch":
			if v.String() != "" {
				method = strings.ToUpper(string(fd.Name()))
				uri = v.String()
			}
		case "custom":
			custom := v.Message()
			custom.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
				switch fd.Name() {
				case "kind":
					method = v.String()
				case "path":
					uri = v.String()
				}
				return true
			})
		}
		return true
	})

	if uri != "" {
		return &HTTPRule{
			Method:       method,
			URI:          uri,
			Body:         body,
			ResponseBody: responseBody,
		}
	}
	return nil
}

// HTTPRule defines a parsed, easier-to-query equivalent to `annotations.HttpRule`.
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
//
// For a given variable, the key is the variable's field path. The value is the
// variable's template, which will match segment(s) of the URL.
//
// For more details on the path template syntax, see
// https://github.com/googleapis/googleapis/blob/6e1a5a066659794f26091674e3668229e7750052/google/api/http.proto#L224.
func (h *HTTPRule) GetVariables() map[string]string {
	vars := map[string]string{}

	// Replace the version template variable with "v".
	uri := VersionedSegment.ReplaceAllString(h.URI, "v")
	for _, match := range plainVar.FindAllStringSubmatch(uri, -1) {
		vars[match[1]] = "*"
	}
	for _, match := range varSegment.FindAllStringSubmatch(uri, -1) {
		vars[match[1]] = match[2]
	}
	return vars
}

// GetPlainURI returns the URI with variable segment information removed.
func (h *HTTPRule) GetPlainURI() string {
	return plainVar.ReplaceAllString(
		varSegment.ReplaceAllString(
			VersionedSegment.ReplaceAllString(h.URI, "v"),
			"$2"),
		"*")
}

var (
	plainVar   = regexp.MustCompile(`\{([^}=]+)\}`)
	varSegment = regexp.MustCompile(`\{([^}=]+)=([^}]+)\}`)
	// VersionedSegment is a regex to extract the API version from
	// an HTTP path.
	VersionedSegment = regexp.MustCompile(`\{\$api_version\}`)
)
