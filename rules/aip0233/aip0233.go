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

// Package aip0233 contains rules defined in https://aip.dev/233.
package aip0233

import (
	"regexp"

	"github.com/googleapis/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
)

// AddRules accepts a register function and registers each of
// this AIP's rules to it.
func AddRules(r lint.RuleRegistry) error {
	return r.Register(
		233,
		pluralMethodName,
		requestMessageName,
		responseMessageName,
		httpBody,
		httpUriSuffix,
		httpVerb,
		requestParentField,
		requestParentReference,
		requestRequestsBehavior,
		requestRequestsField,
		requestUnknownFields,
		resourceReferenceType,
		responseResourceField,
	)
}

var batchCreateMethodRegexp = regexp.MustCompile("^BatchCreate(?:[A-Za-z0-9]|$)")
var batchCreateReqMessageRegexp = regexp.MustCompile("^BatchCreate[A-Za-z0-9]*Request$")
var batchCreateResMessageRegexp = regexp.MustCompile("^BatchCreate[A-Za-z0-9]*Response$")
var batchCreateURINameRegexp = regexp.MustCompile(`:batchCreate$`)

// Returns true if this is a AIP-233 Batch Create method, false otherwise.
func isBatchCreateMethod(m *desc.MethodDescriptor) bool {
	return batchCreateMethodRegexp.MatchString(m.GetName())
}

// Returns true if this is an AIP-233 Batch Create request message, false otherwise.
func isBatchCreateRequestMessage(m *desc.MessageDescriptor) bool {
	return batchCreateReqMessageRegexp.MatchString(m.GetName())
}

// Returns true if this is an AIP-233 Batch Create response message, false otherwise.
func isBatchCreateResponseMessage(m *desc.MessageDescriptor) bool {
	return batchCreateResMessageRegexp.MatchString(m.GetName())
}
