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

	"github.com/jhump/protoreflect/desc"
	"github.com/stoewer/go-strcase"
	apb "google.golang.org/genproto/googleapis/api/annotations"
)

// IsDeclarativeFriendly returns true if the descriptor is declarative-friendly.
//
// For messages, this means the message is annotated with google.api.resource
// and style: DECLARATIVE_FRIENDLY is set.
//
// For methods, the output message is checked and the same logic applied.
// Additionally, methods beginning with "List" (and "Delete" if returning Empty)
// have special logic applied to try to find the correct message and run the
// check.
func IsDeclarativeFriendly(d desc.Descriptor) bool {
	switch m := d.(type) {
	case *desc.MessageDescriptor:
		// We have a message; get the google.api.resource annotation and
		// see if it is styled declarative-friendly.
		if resource := GetResource(m); resource != nil {
			for _, style := range resource.GetStyle() {
				if style == apb.ResourceDescriptor_DECLARATIVE_FRIENDLY {
					return true
				}
			}
		}
	case *desc.MethodDescriptor:
		// Get the return type for the method. If the method is an LRO, then
		// get the response type from the operation_info annotation.
		response := m.GetOutputType()
		if response.GetFullyQualifiedName() == "google.longrunning.Operation" {
			if opInfo := GetOperationInfo(m); opInfo != nil {
				response = findMessage(m.GetFile(), opInfo.GetResponseType())

				// Sanity check: We may not have found the message.
				// If that is the case, give up and assume the method is not
				// declarative-friendly.
				if response == nil {
					return false
				}
			}
		}

		// If the return value has a google.api.resource annotation, we can
		// assume it is the resource and check it.
		if IsResource(response) {
			return IsDeclarativeFriendly(response)
		}

		// If the return value is a List response (AIP-132), we should be able
		// to find the resource as a field in the response.
		if n := response.GetName(); strings.HasPrefix(n, "List") && strings.HasSuffix(n, "Response") {
			for _, field := range response.GetFields() {
				if field.IsRepeated() && field.GetMessageType() != nil {
					return IsDeclarativeFriendly(field.GetMessageType())
				}
			}
		}

		// If this is a Delete method (AIP-135) with a return value of Empty,
		// try to find the resource.
		if strings.HasPrefix(m.GetName(), "Delete") && m.GetOutputType().GetName() == "Empty" {
			if resource := findMessage(m.GetFile(), m.GetName()[6:]); resource != nil {
				return IsDeclarativeFriendly(resource)
			}
		}

		// At this point, we probably have a custom method.
		// Try to identify a resource by whittling away at the method name and
		// seeing if there is a match.
		snakeName := strings.Split(strcase.SnakeCase(m.GetName()), "_")
		for i := 1; i < len(snakeName); i++ {
			name := strcase.UpperCamelCase(strings.Join(snakeName[i:], "_"))
			if resource := findMessage(m.GetFile(), name); resource != nil {
				return IsDeclarativeFriendly(resource)
			}
		}
	}

	// We found no evidence that this descriptor is supposed to be
	// declarative-friendly. Return false.
	return false
}

// findMessage looks for a message in a file and all imports within the
// same package.
func findMessage(f *desc.FileDescriptor, name string) *desc.MessageDescriptor {
	// FileDescriptor.FindMessage requires fully-qualified message names;
	// attempt to infer that.
	if !strings.Contains(name, ".") && f.GetPackage() != "" {
		name = f.GetPackage() + "." + name
	}

	// Attempt to find the message in the file provided.
	if m := f.FindMessage(name); m != nil {
		return m
	}

	// Attempt to find the message in any dependency files if they are in the
	// same package.
	for _, dep := range f.GetDependencies() {
		if f.GetPackage() == dep.GetPackage() {
			if m := findMessage(dep, name); m != nil {
				return m
			}
		}
	}

	// Whelp, no luck. Too bad.
	return nil
}
