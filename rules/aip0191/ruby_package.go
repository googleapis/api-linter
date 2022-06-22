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

var rubyPackage = &lint.FileRule{
	Name: lint.NewRuleName(191, "ruby-package"),
	OnlyIf: func(f *desc.FileDescriptor) bool {
		fops := f.GetFileOptions()
		return fops != nil && fops.GetRubyPackage() != ""
	},
	LintFile: func(f *desc.FileDescriptor) []lint.Problem {
		ns := f.GetFileOptions().GetRubyPackage()
		delim := "::"

		// Check for invalid characters.
		if !rubyValidChars.MatchString(ns) {
			return []lint.Problem{{
				Message:    "Invalid characters: Ruby packages only allow [A-Za-z0-9:].",
				Descriptor: f,
				Location:   locations.FileRubyPackage(f),
			}}
		}

		// Check that upper camel case is used.
		upperCamel := []string{}
		for _, segment := range strings.Split(ns, delim) {
			upperCamel = append(upperCamel, strcase.UpperCamelCase(segment))
		}
		if want := strings.Join(upperCamel, delim); ns != want {
			return []lint.Problem{{
				Message:    "Ruby packages use UpperCamelCase.",
				Suggestion: fmt.Sprintf("option ruby_package = %q;", want),
				Descriptor: f,
				Location:   locations.FileRubyPackage(f),
			}}
		}

		for _, s := range f.GetServices() {
			n := s.GetName()
			if !packagingServiceNameEquals(n, ns, delim) {
				msg := fmt.Sprintf("Casing of Ruby package and service name %q must match.", n)
				return []lint.Problem{{
					Message:    msg,
					Descriptor: f,
					Location:   locations.FileRubyPackage(f),
				}}
			}
		}

		return nil
	},
}

var rubyValidChars = regexp.MustCompile("^[A-Za-z0-9:]+$")
