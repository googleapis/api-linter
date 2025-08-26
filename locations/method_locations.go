// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package locations

import (
	lrpb "cloud.google.com/go/longrunning/autogen/longrunningpb"
	apb "google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/reflect/protoreflect"
	dpb "google.golang.org/protobuf/types/descriptorpb"
)

// MethodRequestType returns the precise location of the method's input type.
func MethodRequestType(m protoreflect.MethodDescriptor) *dpb.SourceCodeInfo_Location {
	return pathLocation(m, 2) // MethodDecriptor.input_type == 2
}

// MethodResponseType returns the precise location of the method's output type.
func MethodResponseType(m protoreflect.MethodDescriptor) *dpb.SourceCodeInfo_Location {
	return pathLocation(m, 3) // MethodDescriptor.output_type == 3
}

// MethodHTTPRule returns the precise location of the method's `google.api.http`
// rule, if any.
func MethodHTTPRule(m protoreflect.MethodDescriptor) *dpb.SourceCodeInfo_Location {
	return MethodOption(m, int(apb.E_Http.TypeDescriptor().Number()))
}

// MethodOperationInfo returns the precise location of the method's
// `google.longrunning.operation_info` annotation, if any.
func MethodOperationInfo(m protoreflect.MethodDescriptor) *dpb.SourceCodeInfo_Location {
	return MethodOption(m, int(lrpb.E_OperationInfo.TypeDescriptor().Number()))
}

// MethodSignature returns the precise location of the method's
// `google.api.method_signature` annotation, if any.
func MethodSignature(m protoreflect.MethodDescriptor, index int) *dpb.SourceCodeInfo_Location {
	return pathLocation(m, 4, int(apb.E_MethodSignature.TypeDescriptor().Number()), index) // MethodDescriptor.options == 4
}

// MethodOption returns the precise location of the method's option with the given field number, if any.
func MethodOption(m protoreflect.MethodDescriptor, fieldNumber int) *dpb.SourceCodeInfo_Location {
	return pathLocation(m, 4, fieldNumber) // MethodDescriptor.options == 4
}
