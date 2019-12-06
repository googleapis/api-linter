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

// Package aip0135 contains rules defined in https://aip.dev/135.
package aip0135

import (
	"regexp"

	"github.com/googleapis/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
)

// AddRules accepts a register function and registers each of
// this AIP's rules to it.
func AddRules(r lint.RuleRegistry) error {
	return r.Register(
		135,
		httpBody,
		httpMethod,
		httpNameField,
		responseMessageName,
		requestMessageName,
		requestNameBehavior,
		standardFields,
		unknownFields,
	)
}

var deleteMethodRegexp = regexp.MustCompile("^Delete(?:[A-Z]|$)")
var deleteReqMessageRegexp = regexp.MustCompile("^Delete[A-Za-z0-9]*Request$")
var deleteURINameRegexp = regexp.MustCompile(`{name=[a-zA-Z/*]+}$`)

// Returns true if this is a AIP-135 Delete method, false otherwise.
func isDeleteMethod(m *desc.MethodDescriptor) bool {
	return deleteMethodRegexp.MatchString(m.GetName())
}

// Returns true if this is an AIP-135 Delete request message, false otherwise.
func isDeleteRequestMessage(m *desc.MessageDescriptor) bool {
	return deleteReqMessageRegexp.MatchString(m.GetName())
}
