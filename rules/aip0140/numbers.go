// Copyright 2020 Google LLC
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
	"regexp"
	"strings"

	"github.com/commure/api-linter/lint"
	"github.com/commure/api-linter/locations"
	"github.com/jhump/protoreflect/desc"
)

var numbers = &lint.FieldRule{
	Name: lint.NewRuleName(140, "numbers"),
	LintField: func(f *desc.FieldDescriptor) []lint.Problem {
		for _, segment := range strings.Split(f.GetName(), "_") {
			if numberStart.MatchString(segment) {
				return []lint.Problem{{
					Message:    "No word in a field name should start with a number.",
					Descriptor: f,
					Location:   locations.DescriptorName(f),
				}}
			}
		}
		return nil
	},
}

var numberStart = regexp.MustCompile("^[0-9]")
