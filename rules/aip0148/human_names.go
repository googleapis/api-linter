// Copyright 2020 Google LLC
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

package aip0148

import (
	"fmt"

	"github.com/commure/api-linter/lint"
	"github.com/commure/api-linter/locations"
	"github.com/jhump/protoreflect/desc"
)

var humanNames = &lint.FieldRule{
	Name: lint.NewRuleName(148, "human-names"),
	LintField: func(f *desc.FieldDescriptor) []lint.Problem {
		for got, want := range corrections {
			if f.GetName() == got {
				return []lint.Problem{{
					Message:    fmt.Sprintf("Use %s instead of %s.", want, got),
					Descriptor: f,
					Location:   locations.DescriptorName(f),
					Suggestion: want,
				}}
			}
		}
		return nil
	},
}

var corrections = map[string]string{
	"first_name": "given_name",
	"last_name":  "family_name",
}
