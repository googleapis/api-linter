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

package aip0140

import (
	"strings"

	"github.com/googleapis/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
)

var base64 = &lint.FieldRule{
	Name:   lint.NewRuleName("core", "0140", "base64"),
	URI:    "https://aip.dev/140",
	OnlyIf: isStringField,
	LintField: func(f *desc.FieldDescriptor) []lint.Problem {
		comment := strings.ToLower(f.GetSourceInfo().GetLeadingComments())
		if strings.Contains(comment, "base64") || strings.Contains(comment, "base-64") {
			return []lint.Problem{{
				Message:    "Field %q mentions base64 encoding in comments, so it should probably be type `bytes`, not `string`.",
				Descriptor: f,
			}}
		}
		return nil
	},
}
