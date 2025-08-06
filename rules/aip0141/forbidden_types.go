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

	"bitbucket.org/creachadair/stringset"
	"github.com/googleapis/api-linter/v2/lint"
	"github.com/googleapis/api-linter/v2/locations"
	"github.com/googleapis/api-linter/v2/rules/internal/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var forbiddenTypes = &lint.FieldRule{
	Name: lint.NewRuleName(141, "forbidden-types"),
	LintField: func(f protoreflect.FieldDescriptor) []lint.Problem {
		nope := stringset.New("fixed32", "fixed64", "uint32", "uint64")
		if typeName := utils.GetTypeName(f); nope.Contains(typeName) {
			// Preserve original intent w/r/t 32-bit vs. 64-bit.
			want := "int" + typeName[len(typeName)-2:]
			return []lint.Problem{{
				Message:    fmt.Sprintf("Use %q instead of %q.", want, typeName),
				Suggestion: want,
				Descriptor: f,
				Location:   locations.FieldType(f),
			}}
		}
		return nil
	},
}
