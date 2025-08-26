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

	"github.com/googleapis/api-linter/v2/lint"
	"github.com/googleapis/api-linter/v2/locations"
	"github.com/stoewer/go-strcase"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"google.golang.org/protobuf/reflect/protoreflect"
	dpb "google.golang.org/protobuf/types/descriptorpb"
)

var csharpNamespace = &lint.FileRule{
	Name: lint.NewRuleName(191, "csharp-namespace"),
	OnlyIf: func(f protoreflect.FileDescriptor) bool {
		return f.Options().(*dpb.FileOptions).GetCsharpNamespace() != ""
	},
	LintFile: func(f protoreflect.FileDescriptor) []lint.Problem {
		ns := f.Options().(*dpb.FileOptions).GetCsharpNamespace()
		delim := "."

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
		for _, segment := range strings.Split(ns, delim) {
			wantSegment := csharpVersionRegexp.ReplaceAllStringFunc(
				strcase.UpperCamelCase(segment),
				func(s string) string {
					point := csharpVersionRegexp.FindStringSubmatch(s)[1]
					if point != "" {
						s = strings.ReplaceAll(s, point, strings.ToUpper(point))
					}
					stability := csharpVersionRegexp.FindStringSubmatch(s)[2]
					title := cases.Title(language.AmericanEnglish)
					return strings.ReplaceAll(s, stability, title.String(stability))
				},
			)
			upperCamel = append(upperCamel, wantSegment)
		}
		if want := strings.Join(upperCamel, delim); ns != want {
			return []lint.Problem{{
				Message:    "C# namespaces use UpperCamelCase.",
				Suggestion: fmt.Sprintf("option csharp_namespace = %q;", want),
				Descriptor: f,
				Location:   locations.FileCsharpNamespace(f),
			}}
		}

		for i := 0; i < f.Services().Len(); i++ {
			s := f.Services().Get(i)
			n := string(s.Name())
			if !packagingServiceNameEquals(n, ns, delim) {
				msg := fmt.Sprintf("Case of C# namespace and service name %q must match.", n)
				return []lint.Problem{{
					Message:    msg,
					Descriptor: f,
					Location:   locations.FileCsharpNamespace(f),
				}}
			}
		}

		return nil
	},
}

var (
	csharpValidChars    = regexp.MustCompile("^[A-Za-z0-9.]+$")
	csharpVersionRegexp = regexp.MustCompile("[0-9]+(p[0-9]+)?(alpha|beta|main)")
)
