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

package aip0191

import (
	"path/filepath"
	"strings"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/jhump/protoreflect/desc"
)

var filename = &lint.FileRule{
	Name: lint.NewRuleName(191, "filenames"),
	LintFile: func(f *desc.FileDescriptor) []lint.Problem {
		fn := strings.ReplaceAll(filepath.Base(f.GetName()), ".proto", "")
		if versionRegexp.MatchString(fn) {
			return []lint.Problem{{
				Message:    "The proto version must not be used as the filename.",
				Descriptor: f,
				Location:   locations.FilePackage(f),
			}}
		}
		if !validCharacterRegexp.MatchString(fn) {
			return []lint.Problem{{
				Message:    "The filename has invalid characters.",
				Descriptor: f,
				Location:   locations.FilePackage(f),
			}}
		}
		return nil
	},
}
