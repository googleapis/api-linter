// Copyright 2020 Google LLC
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

	"github.com/jhump/protoreflect/desc"
)

// IsCommonProto returns true if a proto file is considered "common".
func IsCommonProto(f *desc.FileDescriptor) bool {
	p := f.GetPackage()
	for _, prefix := range []string{"google.api", "google.protobuf", "google.rpc", "google.longrunning"} {
		if strings.HasPrefix(p, prefix) {
			return true
		}
	}
	return false
}
