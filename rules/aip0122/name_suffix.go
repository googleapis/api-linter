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

package aip0122

import (
	"strings"

	"bitbucket.org/creachadair/stringset"
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/jhump/protoreflect/desc"
)

var nameSuffix = &lint.FieldRule{
	Name: lint.NewRuleName(122, "name-suffix"),
	OnlyIf: func(f *desc.FieldDescriptor) bool {
		return !stringset.New("name", "display_name").Contains(f.GetName())
	},
	LintField: func(f *desc.FieldDescriptor) []lint.Problem {
		if n := f.GetName(); strings.HasSuffix(n, "_name") {
			return []lint.Problem{{
				Message:    "Fields should not use the `_name` suffix.",
				Suggestion: strings.TrimSuffix(n, "_name"),
				Descriptor: f,
				Location:   locations.DescriptorName(f),
			}}
		}
		return nil
	},
}
