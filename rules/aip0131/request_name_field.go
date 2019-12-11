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

package aip0131

import (
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"
)

// Get request should have a string name field.
var requestNameField = &lint.FieldRule{
	Name: lint.NewRuleName(131, "request-name-field"),
	OnlyIf: func(f *desc.FieldDescriptor) bool {
		return isGetRequestMessage(f.GetOwner()) && f.GetName() == "name"
	},
	LintField: func(f *desc.FieldDescriptor) []lint.Problem {
		if f.GetType() != builder.FieldTypeString().GetType() {
			return []lint.Problem{{
				Message:    "`name` field on Get RPCs should be a string",
				Descriptor: f,
				Location:   locations.FieldType(f),
				Suggestion: "string",
			}}
		}

		return nil
	},
}
