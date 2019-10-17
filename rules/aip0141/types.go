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

package aip0141

import (
	"fmt"

	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"
)

var forbiddenTypes = &lint.FieldRule{
	Name: lint.NewRuleName("core", "0141", "forbidden-types"),
	URI:  "https://aip.dev/141#guidance",
	LintField: func(f *desc.FieldDescriptor) []lint.Problem {
		// Make a map of the forbidden types.
		nope := make(map[dpb.FieldDescriptorProto_Type]string)
		for _, t := range []*builder.FieldType{
			builder.FieldTypeFixed32(),
			builder.FieldTypeFixed64(),
			builder.FieldTypeUInt32(),
			builder.FieldTypeUInt64(),
		} {
			// Change "TYPE_TYPENAME" to "typename".
			nope[t.GetType()] = utils.GetScalarTypeName(t)
		}
		if typeName, ok := nope[f.GetType()]; ok {
			// Preserve original intent w/r/t 32-bit vs. 64-bit.
			want := "int" + typeName[len(typeName)-2:]
			return []lint.Problem{{
				Message:    fmt.Sprintf("Use %q instead of %q.", want, typeName),
				Descriptor: f,
			}}
		}
		return nil
	},
}
