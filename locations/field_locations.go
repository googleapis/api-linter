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
	apb "google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/reflect/protoreflect"
	dpb "google.golang.org/protobuf/types/descriptorpb"
)

// FieldOption returns the precise location for the given extension definition on
// the given field. This is useful for writing rules against custom extensions.
//
// Example: locations.FieldOption(field, fieldbehaviorpb.E_FieldBehavior)
func FieldOption(f protoreflect.FieldDescriptor, e protoreflect.ExtensionType) *dpb.SourceCodeInfo_Location {
	return pathLocation(f, 8, int(e.TypeDescriptor().Number())) // FieldDescriptor.options == 8
}

// FieldResourceReference returns the precise location for a field's
// resource reference annotation.
func FieldResourceReference(f protoreflect.FieldDescriptor) *dpb.SourceCodeInfo_Location {
	return pathLocation(f, 8, int(apb.E_ResourceReference.TypeDescriptor().Number())) // FieldDescriptor.options == 8
}

// FieldBehavior returns the precise location for a field's
// field_behavior annotation.
func FieldBehavior(f protoreflect.FieldDescriptor) *dpb.SourceCodeInfo_Location {
	return pathLocation(f, 8, int(apb.E_FieldBehavior.TypeDescriptor().Number())) // FieldDescriptor.options == 8
}

// FieldType returns the precise location for a field's type.
func FieldType(f protoreflect.FieldDescriptor) *dpb.SourceCodeInfo_Location {
	if f.Message() != nil || f.Enum() != nil {
		return pathLocation(f, 6) // FieldDescriptor.type_name == 6
	}
	return pathLocation(f, 5) // FieldDescriptor.type == 5
}

// FieldLabel returns the precise location for a field's label.
func FieldLabel(f protoreflect.FieldDescriptor) *dpb.SourceCodeInfo_Location {
	return pathLocation(f, 4) // FieldDescriptor.label == 4
}
