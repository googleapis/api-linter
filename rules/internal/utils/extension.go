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
	apb "google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// GetFieldBehavior returns a stringset.Set of FieldBehavior annotations for
// the given field.
func GetFieldBehavior(f protoreflect.FieldDescriptor) stringset.Set {
	opts := f.Options()
	if !opts.ProtoReflect().Has(apb.E_FieldBehavior.TypeDescriptor()) {
		return stringset.New()
	}
	extValue := opts.ProtoReflect().Get(apb.E_FieldBehavior.TypeDescriptor())
	extList := extValue.List()

	answer := stringset.New()
	for i := 0; i < extList.Len(); i++ {
		fb := apb.FieldBehavior(extList.Get(i).Enum())
		answer.Add(fb.String())
	}
	return answer
}

func getExtensionGeneric[T proto.Message](o protoreflect.Message, ed protoreflect.FieldDescriptor, c T) (T, bool) {
	if !o.Has(ed) {
		var zero T
		return zero, false
	}

	ext := o.Get(ed).Message().Interface()
	if v, ok := ext.(T); ok {
		return v, ok
	}

	d, err := proto.Marshal(ext)
	if err != nil {
		var zero T
		return zero, false
	}
	if err := proto.Unmarshal(d, c); err != nil {
		var zero T
		return zero, false
	}
	return c, true
}

// GetOperationInfo returns the google.longrunning.operation_info annotation.
func GetOperationInfo(m protoreflect.MethodDescriptor) *lrpb.OperationInfo {
	if m == nil {
		return nil
	}
	opInfo, ok := &lrpb.OperationInfo{}, false
	opInfo, ok = getExtensionGeneric(m.Options().ProtoReflect(), lrpb.E_OperationInfo.TypeDescriptor(), opInfo)
	if !ok {
		return nil
	}
	return opInfo
}

// GetOperationResponseType returns the message referred to by the
// (google.longrunning.operation_info).response_type annotation.
func GetOperationResponseType(m protoreflect.MethodDescriptor) protoreflect.MessageDescriptor {
	if m == nil {
		return nil
	}
	info := GetOperationInfo(m)
	if info == nil {
		return nil
	}
	typ := FindMessage(m.ParentFile(), info.GetResponseType())

	return typ
}

// GetResponseType returns the OutputType if the response is
// not an LRO, or the ResponseType otherwise.
func GetResponseType(m protoreflect.MethodDescriptor) protoreflect.MessageDescriptor {
	if m == nil {
		return nil
	}

	ot := m.Output()
	if !isLongRunningOperation(ot) {
		return ot
	}

	return GetOperationResponseType(m)
}

func isLongRunningOperation(m protoreflect.MessageDescriptor) bool {
	return m.ParentFile().Package() == "google.longrunning" && m.Name() == "Operation"
}

// GetMetadataType returns the message referred to by the
// (google.longrunning.operation_info).metadata_type annotation.
func GetMetadataType(m protoreflect.MethodDescriptor) protoreflect.MessageDescriptor {
	if m == nil {
		return nil
	}
	info := GetOperationInfo(m)
	if info == nil {
		return nil
	}
	typ := FindMessage(m.ParentFile(), info.GetMetadataType())

	return typ
}

// GetMethodSignatures returns the `google.api.method_signature` annotations.
func GetMethodSignatures(m protoreflect.MethodDescriptor) [][]string {
	opts := m.Options()
	if !opts.ProtoReflect().Has(apb.E_MethodSignature.TypeDescriptor()) {
		return [][]string{}
	}
	extValue := opts.ProtoReflect().Get(apb.E_MethodSignature.TypeDescriptor())
	extList := extValue.List()

	answer := [][]string{}
	for i := 0; i < extList.Len(); i++ {
		sig := extList.Get(i).String()
		answer = append(answer, strings.Split(sig, ","))
	}
	return answer
}

// GetResource returns the google.api.resource annotation.
func GetResource(m protoreflect.MessageDescriptor) *apb.ResourceDescriptor {
	if m == nil {
		return nil
	}
	res, ok := &apb.ResourceDescriptor{}, false
	if res, ok = getExtensionGeneric(m.Options().ProtoReflect(), apb.E_Resource.TypeDescriptor(), res); !ok {
		return nil
	}

	return res
}

// IsResource returns true if the message has a populated google.api.resource
// annotation with a non-empty "type" field.
func IsResource(m protoreflect.MessageDescriptor) bool {
	if res := GetResource(m); res != nil {
		return res.GetType() != ""
	}
	return false
}

// IsSingletonResource returns true if the given message is a singleton
// resource according to its pattern.
func IsSingletonResource(m protoreflect.MessageDescriptor) bool {
	for _, pattern := range GetResource(m).GetPattern() {
		if IsSingletonResourcePattern(pattern) {
			return true
		}
	}
	return false
}

// IsSingletonResourcePattern returns true if the given message is a singleton
// resource according to its pattern.
func IsSingletonResourcePattern(pattern string) bool {
	// If the pattern ends in something other than "}", that indicates that this is a singleton.
	//
	// For example:
	//   publishers/{publisher}/books/{book} -- not a singleton, many books
	//   publishers/*/settings -- a singleton; one settings object per publisher
	return !strings.HasSuffix(pattern, "}")
}

// GetResourceDefinitions returns the google.api.resource_definition annotations
// for a file.
func GetResourceDefinitions(f protoreflect.FileDescriptor) []*apb.ResourceDescriptor {
	opts := f.Options()
	if !opts.ProtoReflect().Has(apb.E_ResourceDefinition.TypeDescriptor()) {
		return nil
	}
	extValue := opts.ProtoReflect().Get(apb.E_ResourceDefinition.TypeDescriptor())
	extList := extValue.List()

	answer := []*apb.ResourceDescriptor{}
	for i := 0; i < extList.Len(); i++ {
		msg := extList.Get(i).Message().Interface()
		if rd, ok := msg.(*apb.ResourceDescriptor); ok {
			answer = append(answer, rd)
		} else {
			// It may be a dynamic message, so we need to marshal and unmarshal.
			b, err := proto.Marshal(msg)
			if err != nil {
				continue
			}
			rd := &apb.ResourceDescriptor{}
			if err := proto.Unmarshal(b, rd); err != nil {
				continue
			}
			answer = append(answer, rd)
		}
	}
	return answer
}

// HasResourceReference returns if the field has a google.api.resource_reference annotation.
func HasResourceReference(f protoreflect.FieldDescriptor) bool {
	if f == nil {
		return false
	}
	return f.Options().ProtoReflect().Has(apb.E_ResourceReference.TypeDescriptor())
}

// GetResourceReference returns the google.api.resource_reference annotation.
func GetResourceReference(f protoreflect.FieldDescriptor) *apb.ResourceReference {
	if f == nil {
		return nil
	}

	ref, ok := &apb.ResourceReference{}, false
	if ref, ok = getExtensionGeneric(f.Options().ProtoReflect(), apb.E_ResourceReference.TypeDescriptor(), ref); !ok {
		return nil
	}

	return ref
}

// FindResource returns first resource of type matching the reference param.
// resource Type name being referenced. It looks within a given file and its
// depenedencies, it cannot search within the entire protobuf package.
// This is especially useful for resolving google.api.resource_reference
// annotations.
func FindResource(reference string, file protoreflect.FileDescriptor) *apb.ResourceDescriptor {
	m := FindResourceMessage(reference, file)
	return GetResource(m)
}

// FindResourceMessage returns the message containing the first resource of type
// matching the resource Type name being referenced. It looks within a given
// file and its depenedencies, it cannot search within the entire protobuf
// package. This is especially useful for resolving
// google.api.resource_reference annotations to the message that owns a
// resource.
func FindResourceMessage(reference string, file protoreflect.FileDescriptor) protoreflect.MessageDescriptor {
	files := []protoreflect.FileDescriptor{file}
	for i := 0; i < file.Imports().Len(); i++ {
		files = append(files, file.Imports().Get(i).FileDescriptor)
	}

	for _, f := range files {
		for i := 0; i < f.Messages().Len(); i++ {
			m := f.Messages().Get(i)
			if r := GetResource(m); r != nil {
				if r.GetType() == reference {
					return m
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
func FindResourceChildren(parent *apb.ResourceDescriptor, file protoreflect.FileDescriptor) []*apb.ResourceDescriptor {
	pats := parent.GetPattern()
	if len(pats) == 0 {
		return nil
	}
	// Use the first pattern in the resource because:
	// 1. Patterns cannot be rearranged, so this is the true first pattern
	// 2. The true first pattern is the one most likely to be used as a parent.
	first := pats[0]

	var children []*apb.ResourceDescriptor
	files := []protoreflect.FileDescriptor{file}
	for i := 0; i < file.Imports().Len(); i++ {
		files = append(files, file.Imports().Get(i).FileDescriptor)
	}

	for _, f := range files {
		for i := 0; i < f.Messages().Len(); i++ {
			m := f.Messages().Get(i)
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

func HasFieldInfo(fd protoreflect.FieldDescriptor) bool {
	return fd != nil && fd.Options().ProtoReflect().Has(apb.E_FieldInfo.TypeDescriptor())
}

func GetFieldInfo(fd protoreflect.FieldDescriptor) *apb.FieldInfo {
	if !HasFieldInfo(fd) {
		return nil
	}

	fi, ok := &apb.FieldInfo{}, false
	if fi, ok = getExtensionGeneric(fd.Options().ProtoReflect(), apb.E_FieldInfo.TypeDescriptor(), fi); !ok {
		return nil
	}

	return fi
}

func HasFormat(fd protoreflect.FieldDescriptor) bool {
	if !HasFieldInfo(fd) {
		return false
	}

	fi := GetFieldInfo(fd)
	return fi.GetFormat() != apb.FieldInfo_FORMAT_UNSPECIFIED
}

func GetFormat(fd protoreflect.FieldDescriptor) apb.FieldInfo_Format {
	if !HasFormat(fd) {
		return apb.FieldInfo_FORMAT_UNSPECIFIED
	}
	return GetFieldInfo(fd).GetFormat()
}
