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

package aip0143

import (
	"fmt"

	"github.com/commure/api-linter/lint"
	"github.com/commure/api-linter/locations"
	"github.com/jhump/protoreflect/desc"
)

var fieldNames = &lint.FieldRule{
	Name: lint.NewRuleName(143, "standardized-codes"),
	LintField: func(f *desc.FieldDescriptor) []lint.Problem {
		variants := map[string]string{
			"content_type": "mime_type",
			"country":      "region_code",
			"country_code": "region_code",
			"currency":     "currency_code",
			"lang":         "language_code",
			"language":     "language_code",
			"mime":         "mime_type",
			"mimetype":     "mime_type",
			"tz":           "time_zone",
			"timezone":     "time_zone",
		}
		if want, ok := variants[f.GetName()]; ok {
			return []lint.Problem{{
				Message:    fmt.Sprintf("Use %q in place of %q.", want, f.GetName()),
				Descriptor: f,
				Location:   locations.DescriptorName(f),
				Suggestion: want,
			}}
		}
		return nil
	},
}
