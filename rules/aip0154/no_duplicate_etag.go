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

	"github.com/googleapis/api-linter/v2/lint"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var noDuplicateEtag = &lint.FieldRule{
	Name: lint.NewRuleName(154, "no-duplicate-etag"),
	OnlyIf: func(f protoreflect.FieldDescriptor) bool {
		return string(f.Name()) == "etag" && strings.HasSuffix(string(f.Parent().Name()), "Request")
	},
	LintField: func(f protoreflect.FieldDescriptor) []lint.Problem {
		for i := 0; i < f.Parent().(protoreflect.MessageDescriptor).Fields().Len(); i++ {
			otherField := f.Parent().(protoreflect.MessageDescriptor).Fields().Get(i)
			if m := otherField.Message(); m != nil {
				// if strings.Contains("UpdateBookRequest", "Book")
				//
				// If this is a random, unrelated (not the resource) message, we want to ignore it.
				// Ditto for *other* resources, which could be relevant for custom methods,
				// which is why we do a string check and not a google.api.resource check.
				if strings.Contains(string(f.Parent().Name()), string(m.Name())) && m.Fields().ByName("etag") != nil {
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
