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
func GetHTTPRules(m *desc.MethodDescriptor) []*apb.HttpRule {
	rules := []*apb.HttpRule{}

	// Get the method options.
	opts := m.GetMethodOptions()

	// Get the "primary" rule (the direct google.api.http annotation).
	if x, err := proto.GetExtension(opts, apb.E_Http); err == nil {
		httpRule := x.(*apb.HttpRule)
		rules = append(rules, httpRule)

		// Add any additional bindings and flatten them into `rules`.
		rules = append(rules, httpRule.GetAdditionalBindings()...)
	}

	// Done; return the rules.
	return rules
}
