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
	"github.com/googleapis/api-linter/v2/lint"
	"github.com/googleapis/api-linter/v2/locations"
	"github.com/googleapis/api-linter/v2/rules/internal/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var fieldTypes = &lint.FieldRule{
	Name: lint.NewRuleName(143, "string-type"),
	OnlyIf: func(f protoreflect.FieldDescriptor) bool {
		return stringset.New(
			"country_code",
			"currency_code",
			"language_code",
			"mime_type",
			"time_zone",
		).Contains(string(f.Name()))
	},
	LintField: func(f protoreflect.FieldDescriptor) []lint.Problem {
		if typeName := utils.GetTypeName(f); typeName != "string" {
			if f.Name() == "time_zone" && typeName == "google.type.TimeZone" {
				return nil
			}
			return []lint.Problem{{
				Message:    fmt.Sprintf("Field %q should be a string, not %s.", f.Name(), typeName),
				Suggestion: "string",
				Descriptor: f,
				Location:   locations.FieldType(f),
			}}
		}
		return nil
	},
}
