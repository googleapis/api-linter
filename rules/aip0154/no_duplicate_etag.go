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

package aip0154

import (
	"strings"

	"github.com/googleapis/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
)

var noDuplicateEtag = &lint.FieldRule{
	Name: lint.NewRuleName(154, "no-duplicate-etag"),
	OnlyIf: func(f *desc.FieldDescriptor) bool {
		return f.GetName() == "etag" && strings.HasSuffix(f.GetOwner().GetName(), "Request")
	},
	LintField: func(f *desc.FieldDescriptor) []lint.Problem {
		for _, otherField := range f.GetOwner().GetFields() {
			if m := otherField.GetMessageType(); m != nil {
				// if strings.Contains("UpdateBookRequest", "Book")
				//
				// If this is a random, unrelated (not the resource) message, we want to ignore it.
				// Ditto for *other* resources, which could be relevant for custom methods,
				// which is why we do a string check and not a google.api.resource check.
				if strings.Contains(f.GetOwner().GetName(), m.GetName()) && m.FindFieldByName("etag") != nil {
					return []lint.Problem{{
						Message:    "Request messages that include the resource should omit etag.",
						Descriptor: f,
						Suggestion: "",
					}}
				}
			}
		}
		return nil
	},
}
