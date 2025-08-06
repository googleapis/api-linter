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
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

var rubyPackage = &lint.FileRule{
	Name: lint.NewRuleName(191, "ruby-package"),
	OnlyIf: func(f protoreflect.FileDescriptor) bool {
		return f.Options().(*descriptorpb.FileOptions).GetRubyPackage() != ""
	},
	LintFile: func(f protoreflect.FileDescriptor) []lint.Problem {
		ns := f.Options().(*descriptorpb.FileOptions).GetRubyPackage()
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
			if maybeVersionRegexp.MatchString(segment) {
				// The version portion of the package for Ruby should only have
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
				Message:    "Ruby packages use UpperCamelCase.",
				Suggestion: fmt.Sprintf("option ruby_package = %q;", want),
				Descriptor: f,
				Location:   locations.FileRubyPackage(f),
			}}
		}

		for i := 0; i < f.Services().Len(); i++ {
			s := f.Services().Get(i)
			n := string(s.Name())
			if !packagingServiceNameEquals(n, ns, delim) {
				msg := fmt.Sprintf("Case of Ruby package and service name %q must match.", n)
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
