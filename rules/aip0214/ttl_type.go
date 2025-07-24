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

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var ttlType = &lint.FieldRule{
	Name: lint.NewRuleName(214, "ttl-type"),
	OnlyIf: func(f protoreflect.FieldDescriptor) bool {
		return f.Name() == "ttl"
	},
	LintField: func(f protoreflect.FieldDescriptor) []lint.Problem {
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
