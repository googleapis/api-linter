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

// Package aip0133 contains rules defined in https://aip.dev/133.
package aip0133

import (
	"regexp"
	"strings"

	"github.com/googleapis/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
)

// AddRules accepts a register function and registers each of
// this AIP's rules to it.
func AddRules(r lint.RuleRegistry) {
	r.Register(
		httpUriField,
		httpVerb,
		httpBody,
		inputName,
		outputName,
		resourceField,
		parentField,
		unknownFields,
	)
}

var createMethodRegexp = regexp.MustCompile("^Create(?:[A-Z]|$)")
var createReqMessageRegexp = regexp.MustCompile("^Create[A-Za-z0-9]*Request$")
var createURINameRegexp = regexp.MustCompile("\\{parent=[a-zA-Z/*]+\\}$")

// Returns true if this is a AIP-133 Create method, false otherwise.
func isCreateMethod(m *desc.MethodDescriptor) bool {
	return createMethodRegexp.MatchString(m.GetName())
}

// Returns true if this is an AIP-133 Get request message, false otherwise.
func isCreateRequestMessage(m *desc.MessageDescriptor) bool {
	return createReqMessageRegexp.MatchString(m.GetName())
}

// get resource message type name from method
func getResourceMsgName(m *desc.MethodDescriptor) string {
	if !isCreateMethod(m) {
		return ""
	}

	// Usually the response message will be the resource message, and its name will
	// be part of method name (make a double check here to avoid the issue when
	// method or output naming doesn't follow the right principles)
	if strings.Contains(m.GetName()[6:], m.GetOutputType().GetName()) {
		return m.GetOutputType().GetName()
	} else {
		return m.GetName()[6:]
	}
}

// get resource message type name from request message
func getResourceMsgNameFromReq(m *desc.MessageDescriptor) string {
	if !isCreateRequestMessage(m) {
		return ""
	}

	// retrieve the string between the prefix "Create" and suffix "Request" from
	// the name "Create<XXX>Request", and this part will usually be the resource
	// message name(if its naming follows the right principle)
	resourceMsgName := m.GetName()[6 : len(m.GetName())-7]

	// Get the resource field of the request message if it exist, this part will
	// be exactly the resource message name (make a double check here to avoid the
	// issues when request message naming doesn't follow the right principles)
	for _, fieldDesc := range m.GetFields() {
		if msgDesc := fieldDesc.GetMessageType(); msgDesc != nil && strings.Contains(resourceMsgName, msgDesc.GetName()) {
			resourceMsgName = msgDesc.GetName()
		}
	}

	return resourceMsgName
}
