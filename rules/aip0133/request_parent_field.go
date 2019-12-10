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

package aip0133

import (
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"
)

// The type of the parent field in a creat request should be string.
var requestParentField = &lint.FieldRule{
	Name: lint.NewRuleName(133, "request-parent-field"),
	OnlyIf: func(f *desc.FieldDescriptor) bool {
		return isCreateRequestMessage(f.GetOwner()) && f.GetName() == "parent"
	},
	LintField: func(f *desc.FieldDescriptor) []lint.Problem {
		if f.GetType() != builder.FieldTypeString().GetType() {
			return []lint.Problem{{
				Message:    "`parent` field on create request message must be a string",
				Suggestion: "string",
				Descriptor: f,
				Location:   locations.FieldType(f),
			}}
		}

		return nil
	},
}
