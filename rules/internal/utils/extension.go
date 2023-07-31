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

package utils

import (
	"strings"

	"bitbucket.org/creachadair/stringset"
	lrpb "cloud.google.com/go/longrunning/autogen/longrunningpb"
	"github.com/jhump/protoreflect/desc"
	apb "google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
)

// GetFieldBehavior returns a stringset.Set of FieldBehavior annotations for
// the given field.
func GetFieldBehavior(f *desc.FieldDescriptor) stringset.Set {
	opts := f.GetFieldOptions()
	if x := proto.GetExtension(opts, apb.E_FieldBehavior); x != nil {
		answer := stringset.New()
		for _, fb := range x.([]apb.FieldBehavior) {
			answer.Add(fb.String())
		}
		return answer
	}
	return nil
}

// GetOperationInfo returns the google.longrunning.operation_info annotation.
func GetOperationInfo(m *desc.MethodDescriptor) *lrpb.OperationInfo {
	if m == nil {
		return nil
	}
	opts := m.GetMethodOptions()
	if x := proto.GetExtension(opts, lrpb.E_OperationInfo); x != nil {
		return x.(*lrpb.OperationInfo)
	}
	return nil
}

// GetOperationResponseType returns the message referred to by the
// (google.longrunning.operation_info).response_type annotation.
func GetOperationResponseType(m *desc.MethodDescriptor) *desc.MessageDescriptor {
	if m == nil {
		return nil
	}
	info := GetOperationInfo(m)
	if info == nil {
		return nil
	}
	typ := FindMessage(m.GetFile(), info.GetResponseType())

	return typ
}

// GetResponseType returns the OutputType if the response is
// not an LRO, or the ResponseType otherwise.
func GetResponseType(m *desc.MethodDescriptor) *desc.MessageDescriptor {
	if m == nil {
		return nil
	}

	ot := m.GetOutputType()
	if !isLongRunningOperation(ot) {
		return ot
	}

	return GetOperationResponseType(m)
}

func isLongRunningOperation(m *desc.MessageDescriptor) bool {
	return m.GetFile().GetPackage() == "google.longrunning" && m.GetName() == "Operation"
}

// GetMetadataType returns the message referred to by the
// (google.longrunning.operation_info).metadata_type annotation.
func GetMetadataType(m *desc.MethodDescriptor) *desc.MessageDescriptor {
	if m == nil {
		return nil
	}
	info := GetOperationInfo(m)
	if info == nil {
		return nil
	}
	typ := FindMessage(m.GetFile(), info.GetMetadataType())

	return typ
}

// GetMethodSignatures returns the `google.api.method_signature` annotations.
func GetMethodSignatures(m *desc.MethodDescriptor) [][]string {
	answer := [][]string{}
	opts := m.GetMethodOptions()
	if x := proto.GetExtension(opts, apb.E_MethodSignature); x != nil {
		for _, sig := range x.([]string) {
			answer = append(answer, strings.Split(sig, ","))
		}
	}
	return answer
}

// GetResource returns the google.api.resource annotation.
func GetResource(m *desc.MessageDescriptor) *apb.ResourceDescriptor {
	if m == nil {
		return nil
	}
	opts := m.GetMessageOptions()
	if x := proto.GetExtension(opts, apb.E_Resource); x != nil {
		return x.(*apb.ResourceDescriptor)
	}
	return nil
}

// IsResource returns true if the message has a populated google.api.resource
// annotation with a non-empty "type" field.
func IsResource(m *desc.MessageDescriptor) bool {
	if res := GetResource(m); res != nil {
		return res.GetType() != ""
	}
	return false
}

// IsSingletonResource returns true if the given message is a singleton
// resource according to its pattern.
func IsSingletonResource(m *desc.MessageDescriptor) bool {
	// If the pattern ends in something other than "}", that indicates that this is a singleton.
	//
	// For example:
	//   publishers/{publisher}/books/{book} -- not a singleton, many books
	//   publishers/*/settings -- a singleton; one settings object per publisher
	for _, pattern := range GetResource(m).GetPattern() {
		if !strings.HasSuffix(pattern, "}") {
			return true
		}
	}
	return false
}

// GetResourceDefinitions returns the google.api.resource_definition annotations
// for a file.
func GetResourceDefinitions(f *desc.FileDescriptor) []*apb.ResourceDescriptor {
	opts := f.GetFileOptions()
	if x := proto.GetExtension(opts, apb.E_ResourceDefinition); x != nil {
		return x.([]*apb.ResourceDescriptor)
	}
	return nil
}

// GetResourceReference returns the google.api.resource_reference annotation.
func GetResourceReference(f *desc.FieldDescriptor) *apb.ResourceReference {
	if f == nil {
		return nil
	}
	opts := f.GetFieldOptions()
	if x := proto.GetExtension(opts, apb.E_ResourceReference); x != nil {
		return x.(*apb.ResourceReference)
	}
	return nil
}

// FindResource returns first resource of type matching the reference param.
// resource Type name being referenced. It looks within a given file and its
// depenedencies, it cannot search within the entire protobuf package.
// This is especially useful for resolving google.api.resource_reference
// annotations.
func FindResource(reference string, file *desc.FileDescriptor) *apb.ResourceDescriptor {
	files := append(file.GetDependencies(), file)
	for _, f := range files {
		for _, m := range f.GetMessageTypes() {
			if r := GetResource(m); r != nil {
				if r.GetType() == reference {
					return r
				}
			}
		}
	}
	return nil
}

// SplitResourceTypeName splits the `Resource.type` field into the service name
// and the resource type name.
func SplitResourceTypeName(typ string) (service string, typeName string, ok bool) {
	split := strings.Split(typ, "/")
	if len(split) != 2 || split[0] == "" || split[1] == "" {
		return
	}

	service = split[0]
	typeName = split[1]
	ok = true

	return
}

// FindResourceChildren attempts to search for other resources defined in the
// package that are parented by the given resource.
func FindResourceChildren(parent *apb.ResourceDescriptor, file *desc.FileDescriptor) []*apb.ResourceDescriptor {
	pats := parent.GetPattern()
	if len(pats) == 0 {
		return nil
	}
	// Use the first pattern in the resource because:
	// 1. Patterns cannot be rearranged, so this is the true first pattern
	// 2. The true first pattern is the one most likely to be used as a parent.
	first := pats[0]

	var children []*apb.ResourceDescriptor
	files := append(file.GetDependencies(), file)
	for _, f := range files {
		for _, m := range f.GetMessageTypes() {
			if r := GetResource(m); r != nil && r.GetType() != parent.GetType() {
				for _, p := range r.GetPattern() {
					if strings.HasPrefix(p, first) {
						children = append(children, r)
						break
					}
				}
			}
		}
	}

	return children
}
