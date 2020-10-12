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

// FieldResourceReference returns the precise location for a field's
// resource reference annotation.
func FieldResourceReference(f *desc.FieldDescriptor) *dpb.SourceCodeInfo_Location {
	return pathLocation(f, 8, int(apb.E_ResourceReference.TypeDescriptor().Number())) // FieldDescriptor.options == 8
}

// FieldType returns the precise location for a field's type.
func FieldType(f *desc.FieldDescriptor) *dpb.SourceCodeInfo_Location {
	if f.GetMessageType() != nil || f.GetEnumType() != nil {
		return pathLocation(f, 6) // FieldDescriptor.type_name == 6
	}
	return pathLocation(f, 5) // FieldDescriptor.type == 5
}

// FieldLabel returns the precise location for a field's label.
func FieldLabel(f *desc.FieldDescriptor) *dpb.SourceCodeInfo_Location {
	return pathLocation(f, 4) // FieldDescriptor.label == 4
}
