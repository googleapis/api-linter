// Copyright 2020 Google LLC
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

// Package aip0165 contains rules defined in https://aip.dev/165.
package aip0165

import (
	"regexp"

	"github.com/googleapis/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
)

// AddRules adds all of the AIP-165 rules to the provided registry.
func AddRules(r lint.RuleRegistry) error {
	return r.Register(
		165,
		httpBody,
		httpMethod,
		httpParentVariable,
		httpURISuffix,
		requestFilterBehavior,
		requestFilterField,
		requestForceField,
		requestMessageName,
		requestParentBehavior,
		requestParentField,
		requestParentReference,
		responseMessageName,
		responsePurgeCountField,
		responsePurgeSampleField,
		responsePurgeSampleReference,
	)
}

var (
	purgeMethodRegexp      = regexp.MustCompile("^Purge(?:[A-Z]|$)")
	purgeReqMessageRegexp  = regexp.MustCompile("^Purge[A-Za-z0-9]*Request$")
	purgeRespMessageRegexp = regexp.MustCompile("^Purge[A-Za-z0-9]*Response$")
	purgeURINameRegexp     = regexp.MustCompile(`:purge$`)
)

// Returns true if this is a AIP-165 Purge method, false otherwise.
func isPurgeMethod(m *desc.MethodDescriptor) bool {
	return purgeMethodRegexp.MatchString(m.GetName())
}

// Returns true if this is an AIP-165 Purge request message, false otherwise.
func isPurgeRequestMessage(m *desc.MessageDescriptor) bool {
	return purgeReqMessageRegexp.MatchString(m.GetName())
}

// Returns true if this is an AIP-165 Purge response message, false otherwise.
func isPurgeResponseMessage(m *desc.MessageDescriptor) bool {
	return purgeRespMessageRegexp.MatchString(m.GetName())
}
