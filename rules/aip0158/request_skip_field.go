// Copyright 2021 Google LLC
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

package aip0158

import (
	"github.com/googleapis/api-linter/v2/lint"
	"github.com/googleapis/api-linter/v2/locations"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var requestSkipField = &lint.FieldRule{
	Name: lint.NewRuleName(158, "request-skip-field"),
	OnlyIf: func(f protoreflect.FieldDescriptor) bool {
		return isPaginatedRequestMessage(f.Parent().(protoreflect.MessageDescriptor)) && f.Name() == "skip"
	},
	LintField: func(f protoreflect.FieldDescriptor) (problems []lint.Problem) {
		// Rule check: Ensure that the name page_size is the correct type.
		if f.Kind() != protoreflect.Int32Kind || f.IsList() {
			return []lint.Problem{{
				Message:    "`skip` field on List RPCs should be a singular int32",
				Suggestion: "int32",
				Descriptor: f,
				Location:   locations.FieldType(f),
			}}
		}

		return nil
	},
}
