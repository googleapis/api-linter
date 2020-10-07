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

// Package aip0231 contains rules defined in https://aip.dev/231.
package aip0231

import (
	"regexp"

	"github.com/googleapis/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
)

// AddRules accepts a register function and registers each of
// this AIP's rules to it.
func AddRules(r lint.RuleRegistry) error {
	return r.Register(
		231,
		pluralMethodResourceName,
		inputName,
		outputName,
		httpBody,
		httpVerb,
		namesField,
		requestNamesBehavior,
		requestNamesReference,
		requestParentField,
		requestParentReference,
		requestUnknownFields,
		resourceField,
		uriSuffix,
	)
}

var batchGetMethodRegexp = regexp.MustCompile("^BatchGet(?:[A-Za-z0-9]|$)")
var batchGetReqMessageRegexp = regexp.MustCompile("^BatchGet[A-Za-z0-9]*Request$")
var batchGetResMessageRegexp = regexp.MustCompile("^BatchGet[A-Za-z0-9]*Response$")
var batchGetURINameRegexp = regexp.MustCompile(`:batchGet$`)

// Returns true if this is a AIP-231 Get method, false otherwise.
func isBatchGetMethod(m *desc.MethodDescriptor) bool {
	return batchGetMethodRegexp.MatchString(m.GetName())
}

// Returns true if this is an AIP-231 Get request message, false otherwise.
func isBatchGetRequestMessage(m *desc.MessageDescriptor) bool {
	return batchGetReqMessageRegexp.MatchString(m.GetName())
}

// Returns true if this is an AIP-231 Get request message, false otherwise.
func isBatchGetResponseMessage(m *desc.MessageDescriptor) bool {
	return batchGetResMessageRegexp.MatchString(m.GetName())
}
