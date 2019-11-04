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

// Package aip0140 contains rules defined in https://aip.dev/140.
package aip0140

import (
	"strings"

	"github.com/googleapis/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"
	"github.com/stoewer/go-strcase"
)

// AddRules adds all of the AIP-140 rules to the provided registry.
func AddRules(r lint.RuleRegistry) {
	r.Register(
		abbreviations,
		base64,
		lowerSnake,
		noPrepositions,
	)
}

// toLowerSnakeCase converts s to lower_snake_case.
func toLowerSnakeCase(s string) string {
	return strings.ToLower(strcase.SnakeCase(s))
}

// isStringField returns true if the field is a string field.
func isStringField(f *desc.FieldDescriptor) bool {
	return f.GetType() == builder.FieldTypeString().GetType()
}
