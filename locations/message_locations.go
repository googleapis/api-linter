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
	"github.com/jhump/protoreflect/desc"
	apb "google.golang.org/genproto/googleapis/api/annotations"
	dpb "google.golang.org/protobuf/types/descriptorpb"
)

// MessageResource returns the precise location of the `google.api.resource`
// annotation.
func MessageResource(m *desc.MessageDescriptor) *dpb.SourceCodeInfo_Location {
	return pathLocation(m, 7, int(apb.E_Resource.TypeDescriptor().Number())) // MessageDescriptor.options == 7
}
