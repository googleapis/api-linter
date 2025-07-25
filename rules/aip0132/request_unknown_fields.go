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

package aip0132

import (
	"bitbucket.org/creachadair/stringset"
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var allowedFields = stringset.New(
	"parent",                 // AIP-132
	"page_size",              // AIP-158
	"page_token",             // AIP-158
	"skip",                   // AIP-158
	"filter",                 // AIP-132
	"order_by",               // AIP-132
	"show_deleted",           // AIP-135
	"request_id",             // AIP-155
	"read_mask",              // AIP-157
	"view",                   // AIP-157
	"return_partial_success", // AIP-217
)

// List methods should not have unrecognized fields.
var unknownFields = &lint.FieldRule{
	Name: lint.NewRuleName(132, "request-unknown-fields"),
	OnlyIf: func(f protoreflect.FieldDescriptor) bool {
		if m, ok := f.Parent().(protoreflect.MessageDescriptor); ok {
			return utils.IsListRequestMessage(m)
		}
		return false
	},
	LintField: func(field protoreflect.FieldDescriptor) []lint.Problem {
		if !allowedFields.Contains(string(field.Name())) {
			return []lint.Problem{{
				Message:    "List RPCs should only contain fields explicitly described in AIPs.",
				Descriptor: field,
			}}
		}

		return nil
	},
}