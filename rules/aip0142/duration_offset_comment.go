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
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

var durationOffsetComment = &lint.FieldRule{
	Name: lint.NewRuleName(142, "duration-offset-comment"),
	OnlyIf: func(f *desc.FieldDescriptor) bool {
		return utils.GetTypeName(f) == "google.protobuf.Duration" && strings.HasSuffix(f.GetName(), "_offset")
	},
	LintField: func(f *desc.FieldDescriptor) []lint.Problem {
		comment := strings.ToLower(f.GetSourceInfo().GetLeadingComments())
		if comment == "" || (!strings.Contains(comment, "relative") && !strings.Contains(comment, "in respect") && !strings.Contains(comment, "of the")) {
			return []lint.Problem{{
				Message:    "Duration fields ending in `_offset` must include a clear comment explaining the relative start point.",
				Descriptor: f,
			}}
		}
		return nil
	},
}
