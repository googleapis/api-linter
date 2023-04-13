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

package utils

import (
	"strings"

	"bitbucket.org/creachadair/stringset"
	"github.com/jhump/protoreflect/desc"
	"github.com/stoewer/go-strcase"
	apb "google.golang.org/genproto/googleapis/api/annotations"
)

// DeclarativeFriendlyResource returns the declarative-friendly resource
// associated with this descriptor.
//
// For messages:
// If the message is annotated with google.api.resource and
// style: DECLARATIVE_FRIENDLY is set, that message is returned.
// If the message is a standard method request message for a resource with
// google.api.resource and style:DECLARATIVE_FRIENDLY set, then the resource
// is returned.
//
// For methods:
// If the output message is a declarative-friendly resource, it is returned.
// If the method begins with "List" and the first repeated field is a
// declarative-friendly resource, the resource is returned.
// If the method begins with "Delete", the return type is Empty, and an
// appropriate resource message is found and is declarative-friendly, that
// resource is returned.
// If the method is a custom method where a matching resource is found (by
// subset checks on the name) and is declarative-friendly, the resource is
// returned.
//
// If there is no declarative-friendly resource, it returns nil.
func DeclarativeFriendlyResource(d desc.Descriptor) *desc.MessageDescriptor {
	switch m := d.(type) {
	case *desc.MessageDescriptor:
		// Get the google.api.resource annotation and see if it is styled
		// declarative-friendly.
		if resource := GetResource(m); resource != nil {
			for _, style := range resource.GetStyle() {
				if style == apb.ResourceDescriptor_DECLARATIVE_FRIENDLY {
					return m
				}
			}
		}

		// If this is a standard method request message, find the corresponding
		// resource message. The easiest way to do this is to farm it out to the
		// corresponding method.
		if n := m.GetName(); strings.HasSuffix(n, "Request") {
			if method := FindMethod(m.GetFile(), strings.TrimSuffix(n, "Request")); method != nil {
				return DeclarativeFriendlyResource(method)
			}
		}
	case *desc.MethodDescriptor:
		response := m.GetOutputType()

		// If this is a Delete method (AIP-135) with a return value of Empty,
		// try to find the resource.
		//
		// Note: This needs to precede the LRO logic because Delete requests
		// may resolve to Empty, in which case FindMessage will return nil and
		// short-circuit this logic.
		if strings.HasPrefix(m.GetName(), "Delete") && stringset.New("Empty", "Operation").Contains(m.GetOutputType().GetName()) {
			if resource := FindMessage(m.GetFile(), strings.TrimPrefix(m.GetName(), "Delete")); resource != nil {
				return DeclarativeFriendlyResource(resource)
			}
		}

		// If the method is an LRO, then get the response type from the
		// operation_info annotation.
		if IsOperation(response) {
			if opInfo := GetOperationInfo(m); opInfo != nil {
				response = FindMessage(m.GetFile(), opInfo.GetResponseType())

				// Sanity check: We may not have found the message.
				// If that is the case, give up and assume the method is not
				// declarative-friendly.
				if response == nil {
					return nil
				}
			}
		}

		// If the return value has a google.api.resource annotation, we can
		// assume it is the resource and check it.
		if IsResource(response) {
			return DeclarativeFriendlyResource(response)
		}

		// If the return value is a List response (AIP-132), we should be able
		// to find the resource as a field in the response.
		if n := response.GetName(); strings.HasPrefix(n, "List") && strings.HasSuffix(n, "Response") {
			for _, field := range response.GetFields() {
				if field.IsRepeated() && field.GetMessageType() != nil {
					return DeclarativeFriendlyResource(field.GetMessageType())
				}
			}
		}

		// At this point, we probably have a custom method.
		// Try to identify a resource by whittling away at the method name and
		// seeing if there is a match.
		snakeName := strings.Split(strcase.SnakeCase(m.GetName()), "_")
		for i := 1; i < len(snakeName); i++ {
			name := strcase.UpperCamelCase(strings.Join(snakeName[i:], "_"))
			if resource := FindMessage(m.GetFile(), name); resource != nil {
				return DeclarativeFriendlyResource(resource)
			}
		}
	}
	return nil
}

// IsDeclarativeFriendlyMessage returns true if the descriptor is
// declarative-friendly (if DeclarativeFriendlyResource(m) is not nil).
func IsDeclarativeFriendlyMessage(m *desc.MessageDescriptor) bool {
	return DeclarativeFriendlyResource(m) != nil
}

// IsDeclarativeFriendlyMethod returns true if the method is for a
// declarative-friendly resource (if DeclarativeFriendlyResource(m) is not nil).
func IsDeclarativeFriendlyMethod(m *desc.MethodDescriptor) bool {
	return DeclarativeFriendlyResource(m) != nil
}
