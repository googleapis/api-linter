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
	"strings"

	"github.com/jhump/protoreflect/desc"
)

// GetTypeName returns the name of the type of the field, as a string,
// regardless of primitive, message, etc.
//
// TODO: Add support for map types.
func GetTypeName(f *desc.FieldDescriptor) string {
	if m := f.GetMessageType(); m != nil {
		return m.GetFullyQualifiedName()
	}
	if e := f.GetEnumType(); e != nil {
		return e.GetFullyQualifiedName()
	}
	return strings.ToLower(f.GetType().String()[len("TYPE_"):])
}
