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

package aip0143

import (
	"fmt"

	"bitbucket.org/creachadair/stringset"
	"github.com/commure/api-linter/lint"
	"github.com/commure/api-linter/locations"
	"github.com/commure/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

var fieldTypes = &lint.FieldRule{
	Name: lint.NewRuleName(143, "string-type"),
	OnlyIf: func(f *desc.FieldDescriptor) bool {
		return stringset.New(
			"country_code",
			"currency_code",
			"language_code",
			"mime_type",
			"time_zone",
		).Contains(f.GetName())
	},
	LintField: func(f *desc.FieldDescriptor) []lint.Problem {
		if typeName := utils.GetTypeName(f); typeName != "string" {
			if f.GetName() == "time_zone" && typeName == "google.type.TimeZone" {
				return nil
			}
			return []lint.Problem{{
				Message:    fmt.Sprintf("Field %q should be a string, not %s.", f.GetName(), typeName),
				Suggestion: "string",
				Descriptor: f,
				Location:   locations.FieldType(f),
			}}
		}
		return nil
	},
}
