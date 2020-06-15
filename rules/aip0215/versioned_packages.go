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
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

var versionedPackages = &lint.FileRule{
	Name: lint.NewRuleName(215, "versioned-packages"),
	OnlyIf: func(f *desc.FileDescriptor) bool {
		// Common protos are exempt.
		if utils.IsCommonProto(f) {
			return false
		}

		// Ignore this if there is no package.
		p := strings.Split(f.GetPackage(), ".")
		if len(p) == 1 && p[0] == "" {
			return false
		}

		// Exempt anything ending in .type, or .v1master, .v2master, .master, etc.
		if last := p[len(p)-1]; last == "type" || strings.HasSuffix(last, "master") {
			return false
		}

		// Everything else should follow the rule.
		return true
	},
	LintFile: func(f *desc.FileDescriptor) []lint.Problem {
		if !version.MatchString(f.GetPackage()) {
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
