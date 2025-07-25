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

package aip0132

import (
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var knownFields = map[string]func(protoreflect.FieldDescriptor) []lint.Problem{
	"filter":       utils.LintSingularStringField,
	"order_by":     utils.LintSingularStringField,
	"show_deleted": utils.LintSingularBoolField,
}

// List fields should have the correct type.
var requestFieldTypes = &lint.FieldRule{
	Name: lint.NewRuleName(132, "request-field-types"),
	OnlyIf: func(f protoreflect.FieldDescriptor) bool {
		if m, ok := f.Parent().(protoreflect.MessageDescriptor); ok {
			return utils.IsListRequestMessage(m) && knownFields[string(f.Name())] != nil
		}
		return false
	},
	LintField: func(f protoreflect.FieldDescriptor) []lint.Problem {
		return knownFields[string(f.Name())](f)
	},
}