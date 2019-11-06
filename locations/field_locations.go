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

// FieldType returns the precise location for a field's type.
func FieldType(f *desc.FieldDescriptor) *dpb.SourceCodeInfo_Location {
	if sourceInfo := f.GetSourceInfo(); sourceInfo != nil {
		var path []int32
		if f.GetMessageType() != nil || f.GetEnumType() != nil {
			path = append(sourceInfo.Path, 6) // type_name
		} else {
			path = append(sourceInfo.Path, 5) // type
		}
		return pathLocation(f.GetFile(), path)
	}
	return nil
}
