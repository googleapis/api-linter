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
	"google.golang.org/protobuf/reflect/protoreflect"
)

var nameSuffix = &lint.FieldRule{
	Name: lint.NewRuleName(122, "name-suffix"),
	OnlyIf: func(f protoreflect.FieldDescriptor) bool {
		n := string(f.Name())
		// Ignore `{prefix}_display_name` fields as this seems like a reasonable suffix.
		return strings.HasSuffix(n, "_name") && !strings.HasSuffix(n, "_display_name")
	},
	LintField: func(f protoreflect.FieldDescriptor) []lint.Problem {
		allowedNameFields := stringset.New(
			"display_name",
			"family_name",
			"given_name",
			"full_resource_name",
			"crypto_key_name",
			"cmek_key_name",
		)
		if n := string(f.Name()); !allowedNameFields.Contains(n) {
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
