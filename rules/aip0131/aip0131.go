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

// Package aip0131 contains rules defined in https://aip.dev/131.
package aip0131

import (
	"regexp"

	"github.com/googleapis/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
)

// AddRules accepts a register function and registers each of
// this AIP's rules to it.
func AddRules(r lint.RuleRegistry) error {
	return r.Register(
		131,
		httpBody,
		httpMethod,
		httpNameField,
		responseMessageName,
		requestMessageName,
		requestNameReference,
		standardFields,
		synonyms,
		unknownFields,
	)
}

var getMethodRegexp = regexp.MustCompile("^Get(?:[A-Z]|$)")
var getReqMessageRegexp = regexp.MustCompile("^Get[A-Za-z0-9]*Request$")

// Returns true if this is a AIP-131 Get method, false otherwise.
func isGetMethod(m *desc.MethodDescriptor) bool {
	methodName := m.GetName()
	if methodName == "GetIamPolicy" {
		return false
	}
	return getMethodRegexp.MatchString(methodName)
}

// Returns true if this is an AIP-131 Get request message, false otherwise.
func isGetRequestMessage(m *desc.MessageDescriptor) bool {
	return getReqMessageRegexp.MatchString(m.GetName())
}
