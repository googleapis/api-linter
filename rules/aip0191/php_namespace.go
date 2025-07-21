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

var phpNamespace = &lint.FileRule{
	Name: lint.NewRuleName(191, "php-namespace"),
	OnlyIf: func(f *desc.FileDescriptor) bool {
		fops := f.GetFileOptions()
		return fops != nil && fops.GetPhpNamespace() != ""
	},
	LintFile: func(f *desc.FileDescriptor) []lint.Problem {
		ns := f.GetFileOptions().GetPhpNamespace()
		delim := `\`

		// Check for invalid characters.
		if !phpValidChars.MatchString(ns) {
			return []lint.Problem{{
				Message:    `Invalid characters: PHP namespaces only allow [A-Za-z0-9\].`,
				Descriptor: f,
				Location:   locations.FilePhpNamespace(f),
			}}
		}

		// Check that upper camel case is used.
		upperCamel := []string{}
		for _, segment := range strings.Split(ns, delim) {
			if maybeVersionRegexp.MatchString(segment) {
				// The version portion of the package for PHP should only have
				// the first letter, the 'v', be uppercase. The rest are all
				// lowercased.
				upperCamel = append(upperCamel, strings.ToUpper(segment[:1])+strings.ToLower(segment[1:]))
			} else {
				// None versioned portions follow UpperCamel casing.
				upperCamel = append(upperCamel, strcase.UpperCamelCase(segment))
			}
		}
		if want := strings.Join(upperCamel, delim); ns != want {
			return []lint.Problem{{
				Message: "PHP namespaces use UpperCamelCase.",
				Suggestion: fmt.Sprintf(
					"option php_namespace = %s;",
					// Even though the string value is a single backslash, we want
					// to suggest two backslashes, because that is what should be
					// typed into the editor. We use %s to avoid additional escaping
					// of backslashes by Sprintf.
					strings.ReplaceAll(want, delim, `\\`),
				),
				Descriptor: f,
				Location:   locations.FilePhpNamespace(f),
			}}
		}

		for _, s := range f.GetServices() {
			n := s.GetName()
			if !packagingServiceNameEquals(n, ns, delim) {
				msg := fmt.Sprintf("Case of PHP namespace and service name %q must match.", n)
				return []lint.Problem{{
					Message:    msg,
					Descriptor: f,
					Location:   locations.FilePhpNamespace(f),
				}}
			}
		}

		return nil
	},
}

var phpValidChars = regexp.MustCompile(`^[A-Za-z0-9\\]+$`)
