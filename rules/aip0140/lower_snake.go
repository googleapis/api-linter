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

	"github.com/googleapis/api-linter/v2/lint"
	"github.com/googleapis/api-linter/v2/locations"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Field names must be snake case.
var lowerSnake = &lint.FieldRule{
	Name: lint.NewRuleName(140, "lower-snake"),
	LintField: func(f protoreflect.FieldDescriptor) []lint.Problem {
		if got, want := f.Name(), toLowerSnakeCase(string(f.Name())); string(got) != want {
			return []lint.Problem{{
				Message:    fmt.Sprintf("Field `%s` must use lower_snake_case.", got),
				Suggestion: want,
				Descriptor: f,
				Location:   locations.DescriptorName(f),
			}}
		}
		return nil
	},
}
