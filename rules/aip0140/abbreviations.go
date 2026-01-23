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

	"github.com/googleapis/api-linter/v2/lint"
	"github.com/googleapis/api-linter/v2/locations"
	"github.com/stoewer/go-strcase"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var expectedAbbreviations = map[string]string{
	"configuration": "config",
	"identifier":    "id",
	"information":   "info",
	"specification": "spec",
	"statistics":    "stats",
}

var abbreviations = &lint.DescriptorRule{
	Name: lint.NewRuleName(140, "abbreviations"),
	LintDescriptor: func(d protoreflect.Descriptor) (problems []lint.Problem) {
		// Determine the correct case function to use.
		// Most things in protobuf are PascalCase; the two exceptions are
		// fields (snake case) and enum values (UPPER_CAMEL_CASE).
		//
		// We do not need to worry about word separators though, since
		// we are checking for single words only.
		caseFunc := cases.Title(language.AmericanEnglish).String
		switch d.(type) {
		case protoreflect.FieldDescriptor:
			caseFunc = strings.ToLower
		case protoreflect.EnumValueDescriptor:
			caseFunc = strings.ToUpper
		}

		// Iterate over each abbreviation and determine whether the descriptor's
		// name includes the long name.
		for long, short := range expectedAbbreviations {
			for _, segment := range strings.Split(strcase.SnakeCase(string(d.Name())), "_") {
				if segment == long {
					problems = append(problems, lint.Problem{
						Message: fmt.Sprintf(
							"Use the common abbreviation %q instead of %q.",
							caseFunc(short),
							caseFunc(long),
						),
						Suggestion: strings.ReplaceAll(string(d.Name()), caseFunc(long), caseFunc(short)),
						Descriptor: d,
						Location:   locations.DescriptorName(d),
					})
				}
			}
		}
		return
	},
}
