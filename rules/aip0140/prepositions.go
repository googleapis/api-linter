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
	"strings"

	"bitbucket.org/creachadair/stringset"
	"github.com/googleapis/api-linter/v2/lint"
	"github.com/googleapis/api-linter/v2/locations"
	"github.com/googleapis/api-linter/v2/rules/internal/data"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var noPrepositions = &lint.FieldRule{
	Name: lint.NewRuleName(140, "prepositions"),
	OnlyIf: func(f protoreflect.FieldDescriptor) bool {
		return !stringset.New("order_by", "group_by", "hour_of_day", "day_of_week").Contains(string(f.Name()))
	},
	LintField: func(f protoreflect.FieldDescriptor) (problems []lint.Problem) {
		for _, word := range strings.Split(string(f.Name()), "_") {
			if data.Prepositions.Contains(word) {
				problems = append(problems, lint.Problem{
					Message:    fmt.Sprintf("Avoid using %q in field names.", word),
					Descriptor: f,
					Location:   locations.DescriptorName(f),
				})
			}
		}
		return
	},
}
