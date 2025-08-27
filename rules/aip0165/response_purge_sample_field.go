// Copyright 2020 Google LLC
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

package aip0165

import (
	"fmt"

	"github.com/googleapis/api-linter/v2/lint"
	"github.com/googleapis/api-linter/v2/locations"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// The Purge response message should have purge_sample field.
var responsePurgeSampleField = &lint.MessageRule{
	Name:   lint.NewRuleName(165, "response-purge-sample-field"),
	OnlyIf: isPurgeResponseMessage,
	LintMessage: func(m protoreflect.MessageDescriptor) []lint.Problem {
		// Rule check: Establish that a `purge_sample` field is present.
		field := m.Fields().ByName("purge_sample")
		if field == nil {
			return []lint.Problem{{
				Message:    fmt.Sprintf("Message %q has no `purge_sample` field.", m.Name()),
				Descriptor: m,
			}}
		}

		// Rule check: Establish that the purge_sample field is a repeated string.
		if field.Kind() != protoreflect.StringKind || field.Cardinality() != protoreflect.Repeated {
			return []lint.Problem{{
				Message:    "`purge_sample` field on Purge response message must be a repeated string.",
				Descriptor: field,
				Location:   locations.FieldType(field),
				Suggestion: "repeated string",
			}}
		}

		return nil
	},
}
