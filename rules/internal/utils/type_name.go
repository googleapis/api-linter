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
	"fmt"
	"strings"

	"google.golang.org/protobuf/reflect/protoreflect"
)

// GetTypeName returns the name of the type of the field, as a string,
// regardless of primitive, message, etc.
func GetTypeName(f protoreflect.FieldDescriptor) string {
	if f.IsMap() {
		return fmt.Sprintf("map<%s, %s>", GetTypeName(f.MapKey()), GetTypeName(f.MapValue()))
	}
	if m := f.Message(); m != nil {
		return string(m.FullName())
	}
	if e := f.Enum(); e != nil {
		return string(e.FullName())
	}
	return f.Kind().String()
}

// IsOperation returns if the message is a longrunning Operation or not.
func IsOperation(m protoreflect.MessageDescriptor) bool {
	return m.FullName() == "google.longrunning.Operation"
}

// GetResourceMessageName returns the resource message type name from method
func GetResourceMessageName(m protoreflect.MethodDescriptor, expectedVerb string) string {
	if !strings.HasPrefix(string(m.Name()), expectedVerb) {
		return ""
	}

	// Usually the response message will be the resource message, and its name will
	// be part of method name (make a double check here to avoid the issue when
	// method or output naming doesn't follow the right principles)
	// Ignore this rule if the return type is an LRO
	if strings.Contains(string(m.Name()[len(expectedVerb):]), string(m.Output().Name())) && !IsOperation(m.Output()) {
		return string(m.Output().Name())
	}
	return string(m.Name()[len(expectedVerb):])
}
