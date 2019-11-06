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
)

// MethodRequestType returns the precise location of the method's input type.
func MethodRequestType(m *desc.MethodDescriptor) *dpb.SourceCodeInfo_Location {
	if sourceInfo := m.GetSourceInfo(); sourceInfo != nil {
		path := append(sourceInfo.Path, 2) // input_type == 2
		return pathLocation(m.GetFile(), path)
	}
	return nil
}

// MethodResponseType returns the precise location of the method's output type.
func MethodResponseType(m *desc.MethodDescriptor) *dpb.SourceCodeInfo_Location {
	if sourceInfo := m.GetSourceInfo(); sourceInfo != nil {
		path := append(m.GetSourceInfo().Path, 3) // output_type == 3
		return pathLocation(m.GetFile(), path)
	}
	return nil
}
