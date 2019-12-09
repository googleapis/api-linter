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
	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jhump/protoreflect/desc"
	apb "google.golang.org/genproto/googleapis/api/annotations"
)

// MethodRequestType returns the precise location of the method's input type.
func MethodRequestType(m *desc.MethodDescriptor) *dpb.SourceCodeInfo_Location {
	return pathLocation(m, 2) // MethodDecriptor.input_type == 2
}

// MethodResponseType returns the precise location of the method's output type.
func MethodResponseType(m *desc.MethodDescriptor) *dpb.SourceCodeInfo_Location {
	return pathLocation(m, 3) // MethodDescriptor.output_type == 3
}

// MethodHTTPRule returns the precise location of the method's `google.api.http`
// rule, if any.
func MethodHTTPRule(m *desc.MethodDescriptor) *dpb.SourceCodeInfo_Location {
	return pathLocation(m, 4, int(apb.E_Http.Field)) // MethodDescriptor.options == 4
}

// MethodSignature returns the precise location of the method's
// `google.api.method_signature` annotation, if any.
func MethodSignature(m *desc.MethodDescriptor, index int) *dpb.SourceCodeInfo_Location {
	return pathLocation(m, 4, int(apb.E_MethodSignature.Field), index) // MethodDescriptor.options == 4
}
