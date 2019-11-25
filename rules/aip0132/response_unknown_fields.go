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
	"github.com/jhump/protoreflect/desc"
	"github.com/stoewer/go-strcase"
)

// The resource itself is not included here, but also permitted.
// This is covered in code in the rule itself.
var respAllowedFields = stringset.New(
	"next_page_token",       // AIP-158
	"total_size",            // AIP-132
	"unreachable",           // AIP-217
	"unreachable_locations", // Wrong, but a separate AIP-217 rule catches it.
)

var responseUnknownFields = &lint.FieldRule{
	Name: lint.NewRuleName(132, "response-unknown-fields"),
	OnlyIf: func(f *desc.FieldDescriptor) bool {
		return isListResponseMessage(f.GetOwner())
	},
	LintField: func(f *desc.FieldDescriptor) []lint.Problem {
		// A repeated variant of the resource should be permitted.
		resource := strcase.SnakeCase(listRespMessageRegexp.FindStringSubmatch(f.GetOwner().GetName())[1])
		if f.GetName() == resource {
			return nil
		}

		// It is not the resource field; check it against the whitelist.
		if !respAllowedFields.Contains(f.GetName()) {
			return []lint.Problem{{
				Message:    "List responses should only contain fields explicitly described in AIPs.",
				Descriptor: f,
			}}
		}
		return nil
	},
}
