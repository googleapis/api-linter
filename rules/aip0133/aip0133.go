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
	"strings"

	"github.com/googleapis/api-linter/lint"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// AddRules accepts a register function and registers each of
// this AIP's rules to it.
func AddRules(r lint.RuleRegistry) error {
	return r.Register(
		133,
		httpBody,
		httpURIParent,
		httpURIResource,
		httpMethod,
		inputName,
		methodSignature,
		outputName,
		requestIDField,
		requestParentBehavior,
		requestParentField,
		requestParentReference,
		requestParentRequired,
		requestRequiredFields,
		requestResourceBehavior,
		resourceField,
		resourceReferenceType,
		responseLRO,
		synonyms,
		unknownFields,
	)
}

// get resource message type name from request message
func getResourceMsgNameFromReq(m protoreflect.MessageDescriptor) string {
	// retrieve the string between the prefix "Create" and suffix "Request" from
	// the name "Create<XXX>Request", and this part will usually be the resource
	// message name(if its naming follows the right principle)
	resourceMsgName := string(m.Name())[6 : len(m.Name())-7]

	// Get the resource field of the request message if it exist, this part will
	// be exactly the resource message name (make a double check here to avoid the
	// issues when request message naming doesn't follow the right principles)
	for i := 0; i < m.Fields().Len(); i++ {
		fieldDesc := m.Fields().Get(i)
		if msgDesc := fieldDesc.Message(); msgDesc != nil && strings.Contains(resourceMsgName, string(msgDesc.Name())) {
			resourceMsgName = string(msgDesc.Name())
		}
	}

	return resourceMsgName
}