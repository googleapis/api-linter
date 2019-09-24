// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aip0140

import (
	"fmt"

	"github.com/googleapis/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
)

// Field names must be snake case.
var lowerSnake = &lint.FieldRule{
	Name: lint.NewRuleName("core", "0140", "lower-snake"),
	URI:  "https://aip.dev/140#guidance",
	LintField: func(f *desc.FieldDescriptor) lint.Problems {
		if got, want := f.GetName(), toLowerSnakeCase(f.GetName()); got != want {
			return lint.Problems{{
				Message:    fmt.Sprintf("Field `%s` must use lower_snake_case.", got),
				Suggestion: want,
				Descriptor: f,
			}}
		}
		return nil
	},
}
