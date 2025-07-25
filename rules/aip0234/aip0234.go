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
	"google.golang.org/protobuf/reflect/protoreflect"
)

// AddRules accepts a register function and registers each of
// this AIP's rules to it.
func AddRules(r lint.RuleRegistry) error {
	return r.Register(
		234,
		httpBody,
		httpMethod,
		httpURISuffix,
		pluralMethodName,
		requestMessageName,
		requestParentField,
		requestParentReference,
		requestRequestsBehavior,
		requestRequestsField,
		requestUnknownFields,
		responseMessageName,
		responseResourceField,
	)
}

var (
	batchUpdateMethodRegexp     = regexp.MustCompile("^BatchUpdate(?:[A-Za-z0-9]|$)")
	batchUpdateReqMessageRegexp = regexp.MustCompile("^BatchUpdate[A-Za-z0-9]*Request$")
	batchUpdateResMessageRegexp = regexp.MustCompile("^BatchUpdate[A-Za-z0-9]*Response$")
	batchUpdateURINameRegexp    = regexp.MustCompile(`:batchUpdate$`)
)

// Returns true if this is a AIP-234 Batch Update method, false otherwise.
func isBatchUpdateMethod(m protoreflect.MethodDescriptor) bool {
	return batchUpdateMethodRegexp.MatchString(string(m.Name()))
}

// Returns true if this is an AIP-234 Batch Update request message, false otherwise.
func isBatchUpdateRequestMessage(m protoreflect.MessageDescriptor) bool {
	return batchUpdateReqMessageRegexp.MatchString(string(m.Name()))
}

// Returns true if this is an AIP-234 Batch Update request message, false otherwise.
func isBatchUpdateResponseMessage(m protoreflect.MessageDescriptor) bool {
	return batchUpdateResMessageRegexp.MatchString(string(m.Name()))
}
