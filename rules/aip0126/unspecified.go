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

package aip0126

import (
	"fmt"
	"strings"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/jhump/protoreflect/desc"
	"github.com/stoewer/go-strcase"
)

var unspecified = &lint.EnumRule{
	Name: lint.NewRuleName(126, "unspecified"),
	LintEnum: func(e *desc.EnumDescriptor) []lint.Problem {
		firstValue := e.GetValues()[0]
		want := strings.ToUpper(strcase.SnakeCase(e.GetName()) + "_UNSPECIFIED")
		for _, element := range e.GetValues() {
			if element.GetName() == want && element.GetNumber() == 0 {
				return nil
			}
		}
		return []lint.Problem{{
			Message:    fmt.Sprintf("The first enum value should be %q", want),
			Suggestion: want,
			Descriptor: firstValue,
			Location:   locations.DescriptorName(firstValue),
		}}
	},
}
