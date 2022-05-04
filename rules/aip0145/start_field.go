// Copyright 2022 Google LLC
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

package aip0145

import (
	"strings"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/jhump/protoreflect/desc"
)

var startField = &lint.FieldRule{
	Name: lint.NewRuleName(145, "start-field"),
	OnlyIf: func(f *desc.FieldDescriptor) bool {
		return strings.HasPrefix(f.GetName(), "start")
	},
	LintField: func(f *desc.FieldDescriptor) []lint.Problem {
		name := strings.TrimPrefix(f.GetName(), "start")
		if !strings.HasPrefix(name, "_") || len(name) < 2 {
			return []lint.Problem{{
				Message:    "Fields beginning with `start` should be followed by a type using snake case e.g. start_xxx",
				Descriptor: f,
				Location:   locations.FieldType(f),
			}}
		}
		return nil
	},
}
