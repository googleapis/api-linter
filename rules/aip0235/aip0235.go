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

// Package aip0235 contains rules defined in https://aip.dev/235.
package aip0235

import (
	"regexp"

	"github.com/googleapis/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
)

// AddRules accepts a register function and registers each of
// this AIP's rules to it.
func AddRules(r lint.RuleRegistry) error {
	return r.Register(
		235,
		httpBody,
		httpMethod,
		httpUriSuffix,
		pluralMethodName,
		requestMessageName,
		requestNamesBehavior,
		requestNamesReference,
	)
}

var batchDeleteMethodRegexp = regexp.MustCompile("^BatchDelete(?:[A-Z]|$)")
var batchDeleteReqMessageRegexp = regexp.MustCompile("^BatchDelete[A-Za-z0-9]*Request$")
var batchDeleteURIRegexp = regexp.MustCompile(`:batchDelete$`)

// Returns true if this is a AIP-235 Batch Delete method, false otherwise.
func isBatchDeleteMethod(m *desc.MethodDescriptor) bool {
	return batchDeleteMethodRegexp.MatchString(m.GetName())
}

// Returns true if this is an AIP-235 Batch Delete request message, false otherwise.
func isBatchDeleteRequestMessage(m *desc.MessageDescriptor) bool {
	return batchDeleteReqMessageRegexp.MatchString(m.GetName())
}
