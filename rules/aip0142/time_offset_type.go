// Copyright 2025 Google LLC
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

package aip0142

import (
	"strings"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var timeOffsetType = &lint.FieldRule{
	Name: lint.NewRuleName(142, "time-offset-type"),
	LintField: func(f protoreflect.FieldDescriptor) []lint.Problem {
		if utils.GetTypeName(f) != "google.protobuf.Duration" && strings.HasSuffix(string(f.Name()), "_time_offset") {
			return []lint.Problem{{
				Message:    "Fields with names ending in `_time_offset` must be of type `google.protobuf.Duration`.",
				Descriptor: f,
				Suggestion: "google.protobuf.Duration",
				Location:   locations.FieldType(f),
			}}
		}
		return nil
	},
}
