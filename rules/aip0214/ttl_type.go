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

package aip0214

import (
	"fmt"

	"github.com/commure/api-linter/lint"
	"github.com/commure/api-linter/locations"
	"github.com/commure/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

var ttlType = &lint.FieldRule{
	Name: lint.NewRuleName(214, "ttl-type"),
	OnlyIf: func(f *desc.FieldDescriptor) bool {
		return f.GetName() == "ttl"
	},
	LintField: func(f *desc.FieldDescriptor) []lint.Problem {
		if fieldType := utils.GetTypeName(f); fieldType != "google.protobuf.Duration" {
			return []lint.Problem{{
				Message:    fmt.Sprintf("ttl fields should be `google.protobuf.Duration` type, not `%s`.", fieldType),
				Suggestion: "google.protobuf.Duration",
				Descriptor: f,
				Location:   locations.FieldType(f),
			}}
		}
		return nil
	},
}
