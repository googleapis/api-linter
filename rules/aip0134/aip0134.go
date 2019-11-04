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

// Package aip0134 contains rules defined in https://aip.dev/134.
package aip0134

import (
	"regexp"
	"strings"

	"github.com/googleapis/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
	"github.com/stoewer/go-strcase"
)

// AddRules accepts a register function and registers each of
// this AIP's rules to it.
func AddRules(r lint.RuleRegistry) {
	r.Register(
		httpBody,
		httpMethod,
		httpNameField,
		responseMessageName,
		requestMaskField,
		requestMessageName,
		resourceField,
		synonyms,
		unknownFields,
	)
}

var updateMethodRegexp = regexp.MustCompile("^Update(?:[A-Z]|$)")
var updateReqMessageRegexp = regexp.MustCompile("^Update[A-Za-z0-9]*Request$")
var updateURINameRegexp = regexp.MustCompile("\\{([A-Za-z0-9_]+)\\.name=[a-zA-Z/*]+\\}$")

// Returns true if this is a AIP-134 Update method, false otherwise.
func isUpdateMethod(m *desc.MethodDescriptor) bool {
	methodName := m.GetName()
	return updateMethodRegexp.MatchString(methodName)
}

// Returns true if this is an AIP-134 Update request message, false otherwise.
func isUpdateRequestMessage(m *desc.MessageDescriptor) bool {
	return updateReqMessageRegexp.MatchString(m.GetName())
}

func extractResource(reqName string) string {
	// Strips "Update" from the beginning and "Request" from the end.
	return reqName[6 : len(reqName)-7]
}

func fieldNameFromResource(resource string) string {
	return strings.ToLower(strcase.SnakeCase(resource))
}
