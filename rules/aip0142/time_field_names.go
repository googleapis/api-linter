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

package aip0142

import (
	"strings"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

var fieldNames = &lint.FieldRule{
	Name: lint.NewRuleName(142, "time-field-names"),
	LintField: func(f *desc.FieldDescriptor) []lint.Problem {
		// Look for common non-imperative terms.
		mistakes := map[string]string{
			"created":  "create_time",
			"expired":  "expire_time",
			"modified": "update_time",
			"updated":  "update_time",
		}
		for got, want := range mistakes {
			if strings.Contains(f.GetName(), got) {
				return []lint.Problem{{
					Message:    "Use the imperative mood and a `_time` suffix for timestamps.",
					Descriptor: f,
					Location:   locations.DescriptorName(f),
					Suggestion: want,
				}}
			}
		}

		// Look for timestamps that do not end in `_time`.
		if utils.GetTypeName(f) == "google.protobuf.Timestamp" && !strings.HasSuffix(f.GetName(), "_time") {
			return []lint.Problem{{
				Message:    "Timestamp fields should end in `_time`.",
				Descriptor: f,
				Location:   locations.DescriptorName(f),
			}}
		}

		return nil
	},
}
