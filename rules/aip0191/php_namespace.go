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
	LintFile: func(f *desc.FileDescriptor) []lint.Problem {
		if ns := f.GetFileOptions().GetPhpNamespace(); ns != "" {
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
			for _, segment := range strings.Split(ns, `\`) {
				upperCamel = append(upperCamel, strcase.UpperCamelCase(segment))
			}
			if want := strings.Join(upperCamel, `\`); ns != want {
				return []lint.Problem{{
					Message: "PHP namespaces use UpperCamelCase.",
					Suggestion: fmt.Sprintf(
						"option php_namespace = %q;",
						// Even though the string value is a single backslash, we want
						// to suggest two backslashes, because that is what should be
						// typed into the editor.
						strings.ReplaceAll(want, `\`, `\\`),
					),
					Descriptor: f,
					Location:   locations.FilePhpNamespace(f),
				}}
			}
		}
		return nil
	},
}

var phpValidChars = regexp.MustCompile(`^[A-Za-z0-9\\]+$`)
