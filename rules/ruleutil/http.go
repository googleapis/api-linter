// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ruleutil

import (
	"github.com/golang/protobuf/proto"
	"google.golang.org/genproto/googleapis/api/annotations"
	p "google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

// GetHTTPRules returns a array of HTTP rules for a given method descriptor.
//
// Note: This returns an array -- it takes the google.api.http annotation,
// and then flattens the values in `additional_bindings`. This allows
// rule authors to simply range over all of the HTTP rules, since we almost
// always want to apply the same checks to all of them.
func GetHTTPRules(m p.MethodDescriptor) (rules []*annotations.HttpRule) {
	var httpRule *annotations.HttpRule
	opts := m.Options().(*descriptorpb.MethodOptions)
	if x, err := proto.GetExtension(opts, annotations.E_Http); err == nil {
		httpRule = x.(*annotations.HttpRule)
		rules = append(rules, httpRule)
	}
	for _, additionalBinding := range httpRule.GetAdditionalBindings() {
		rules = append(rules, additionalBinding)
	}
	return rules
}