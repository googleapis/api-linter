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

package aip0158

import (
	"github.com/googleapis/api-linter/lint"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var responseRepeatedFirstField = &lint.MessageRule{
	Name: lint.NewRuleName(158, "response-repeated-first-field"),
	OnlyIf: func(m protoreflect.MessageDescriptor) bool {
		return isPaginatedResponseMessage(m) && len(m.Fields()) > 0
	},
	LintMessage: func(m protoreflect.MessageDescriptor) []lint.Problem {
		// Sanity check: Is the first field (positionally) and the field with
		// an ID of 1 actually the same field?
		if m.Fields()[0] != m.FindFieldByNumber(1) {
			return []lint.Problem{{
				Message:    "The first field of paginated RPCs must have a protobuf ID of 1.",
				Descriptor: m.Fields()[0],
			}}
		}

		// Make sure the field is repeated.
		if !m.Fields()[0].IsRepeated() {
			return []lint.Problem{{
				Message:    "The first field of a paginated response should be repeated.",
				Descriptor: m.Fields()[0],
			}}
		}

		return nil
	},
}
