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

	"github.com/jhump/protoreflect/desc"
)

// GetTypeName returns the name of the type of the field, as a string,
// regardless of primitive, message, etc.
func GetTypeName(f *desc.FieldDescriptor) string {
	if k, v := f.GetMapKeyType(), f.GetMapValueType(); k != nil && v != nil {
		return fmt.Sprintf("map<%s, %s>", GetTypeName(k), GetTypeName(v))
	}
	if m := f.GetMessageType(); m != nil {
		return m.GetFullyQualifiedName()
	}
	if e := f.GetEnumType(); e != nil {
		return e.GetFullyQualifiedName()
	}
	return strings.ToLower(f.GetType().String()[len("TYPE_"):])
}

// IsOperation returns if the message is a longrunning Operation or not.
func IsOperation(m *desc.MessageDescriptor) bool {
	return m.GetFullyQualifiedName() == "google.longrunning.Operation"
}

// GetResourceMessageName returns the resource message type name from method
func GetResourceMessageName(m *desc.MethodDescriptor, expectedVerb string) string {
	if !strings.HasPrefix(m.GetName(), expectedVerb) {
		return ""
	}

	// Usually the response message will be the resource message, and its name will
	// be part of method name (make a double check here to avoid the issue when
	// method or output naming doesn't follow the right principles)
	// Ignore this rule if the return type is an LRO
	if strings.Contains(m.GetName()[len(expectedVerb):], m.GetOutputType().GetName()) && !IsOperation(m.GetOutputType()) {
		return m.GetOutputType().GetName()
	}
	return m.GetName()[len(expectedVerb):]
}
