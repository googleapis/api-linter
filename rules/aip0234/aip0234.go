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

// Package aip0234 contains rules defined in https://aip.dev/234.
package aip0234

import (
	"regexp"

	"github.com/googleapis/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
)

// AddRules accepts a register function and registers each of
// this AIP's rules to it.
func AddRules(r lint.RuleRegistry) error {
	return r.Register(
		234,
		httpBody,
		httpMethod,
		httpUriSuffix,
		pluralMethodName,
		requestMessageName,
		requestParentField,
		requestRequestsField,
		responseMessageName,
		responseResourceField,
	)
}

var batchUpdateMethodRegexp = regexp.MustCompile("^BatchUpdate(?:[A-Za-z0-9]|$)")
var batchUpdateReqMessageRegexp = regexp.MustCompile("^BatchUpdate[A-Za-z0-9]*Request$")
var batchUpdateResMessageRegexp = regexp.MustCompile("^BatchUpdate[A-Za-z0-9]*Response$")
var batchUpdateURINameRegexp = regexp.MustCompile(`:batchUpdate$`)

// Returns true if this is a AIP-234 Batch Update method, false otherwise.
func isBatchUpdateMethod(m *desc.MethodDescriptor) bool {
	return batchUpdateMethodRegexp.MatchString(m.GetName())
}

// Returns true if this is an AIP-234 Batch Update request message, false otherwise.
func isBatchUpdateRequestMessage(m *desc.MessageDescriptor) bool {
	return batchUpdateReqMessageRegexp.MatchString(m.GetName())
}

// Returns true if this is an AIP-234 Batch Update request message, false otherwise.
func isBatchUpdateResponseMessage(m *desc.MessageDescriptor) bool {
	return batchUpdateResMessageRegexp.MatchString(m.GetName())
}
