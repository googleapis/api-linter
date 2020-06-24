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
	"github.com/jhump/protoreflect/desc"
	apb "google.golang.org/genproto/googleapis/api/annotations"
	lrpb "google.golang.org/genproto/googleapis/longrunning"
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
	opts := m.GetMethodOptions()
	if x := proto.GetExtension(opts, lrpb.E_OperationInfo); x != nil {
		return x.(*lrpb.OperationInfo)
	}
	return nil
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
	opts := f.GetFieldOptions()
	if x := proto.GetExtension(opts, apb.E_ResourceReference); x != nil {
		return x.(*apb.ResourceReference)
	}
	return nil
}
