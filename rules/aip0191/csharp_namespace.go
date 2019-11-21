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

package aip0191

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/jhump/protoreflect/desc"
	"github.com/stoewer/go-strcase"
)

var csharpNamespace = &lint.FileRule{
	Name: lint.NewRuleName(191, "csharp-namespace"),
	LintFile: func(f *desc.FileDescriptor) []lint.Problem {
		if ns := f.GetFileOptions().GetCsharpNamespace(); ns != "" {
			// Check for invalid characters.
			if !csharpValidChars.MatchString(ns) {
				return []lint.Problem{{
					Message:    "Invalid characters: C# namespaces only allow [A-Za-z0-9.].",
					Descriptor: f,
					Location:   locations.FileCsharpNamespace(f),
				}}
			}

			// Check that upper camel case is used.
			upperCamel := []string{}
			for _, segment := range strings.Split(ns, ".") {
				wantSegment := csharpVersionRegexp.ReplaceAllStringFunc(
					strcase.UpperCamelCase(segment),
					func(s string) string {
						stability := csharpVersionRegexp.FindStringSubmatch(s)[1]
						return strings.ReplaceAll(s, stability, strings.Title(stability))
					},
				)
				upperCamel = append(upperCamel, wantSegment)
			}
			if want := strings.Join(upperCamel, "."); ns != want {
				return []lint.Problem{{
					Message:    "C# namespaces use UpperCamelCase.",
					Suggestion: fmt.Sprintf("option csharp_namespace = %q;", want),
					Descriptor: f,
					Location:   locations.FileCsharpNamespace(f),
				}}
			}
		}
		return nil
	},
}

var csharpValidChars = regexp.MustCompile("^[A-Za-z0-9.]+$")
var csharpVersionRegexp = regexp.MustCompile("[0-9]+(alpha|beta)")
