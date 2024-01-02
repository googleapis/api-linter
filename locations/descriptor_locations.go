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
	dpb "google.golang.org/protobuf/types/descriptorpb"
)

// DescriptorName returns the precise location for a descriptor's name.
//
// This works for any descriptor, regardless of type (message, field, etc.).
func DescriptorName(d desc.Descriptor) *dpb.SourceCodeInfo_Location {
	// All descriptors seem to have `string name = 1`, so this conveniently works.
	return pathLocation(d, 1)
}
