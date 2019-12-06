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

package aip0215

import (
	"regexp"
	"strings"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/jhump/protoreflect/desc"
)

var versionedPackages = &lint.FileRule{
	Name: lint.NewRuleName(215, "versioned-packages"),
	OnlyIf: func(f *desc.FileDescriptor) bool {
		return f.GetPackage() != ""
	},
	LintFile: func(f *desc.FileDescriptor) []lint.Problem {
		if p := f.GetPackage(); !version.MatchString(p) && !strings.HasSuffix(p, ".type") {
			return []lint.Problem{{
				Message:    "API components should be in versioned packages.",
				Descriptor: f,
				Location:   locations.FilePackage(f),
			}}
		}
		return nil
	},
}

var version = regexp.MustCompile(`\.v[\d]+(p[\d]+)?(alpha|beta|eap|test)?[\d]*$`)
