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

package aip0233

import (
	"fmt"

	"bitbucket.org/creachadair/stringset"

	"github.com/googleapis/api-linter/lint"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Batch Create methods should not have unrecognized fields.
var requestUnknownFields = &lint.FieldRule{
	Name: lint.NewRuleName(233, "request-unknown-fields"),
	OnlyIf: func(f protoreflect.FieldDescriptor) bool {
		if m, ok := f.Parent().(protoreflect.MessageDescriptor); ok {
			return isBatchCreateRequestMessage(m)
		}
		return false
	},
	LintField: func(field protoreflect.FieldDescriptor) []lint.Problem {
		allowedFields := stringset.New(
			"parent",        // AIP-233
			"request_id",    // AIP-155
			"requests",      // AIP-233
			"validate_only", // AIP-163
		)
		if !allowedFields.Contains(string(field.Name())) {
			return []lint.Problem{{
				Message: fmt.Sprintf(
					"Unexpected field: Batch Create RPCs must only contain fields explicitly described in https://aip.dev/233, not %q.",
					string(field.Name()),
				),
				Descriptor: field,
			}}
		}

		return nil
	},
}
