// Copyright 2021 Google LLC
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

package aip0152

import (
	"fmt"
	"strings"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// The name of the resource must end with the word "Job".
var requestResourceSuffix = &lint.FieldRule{
	Name: lint.NewRuleName(152, "request-resource-suffix"),
	OnlyIf: func(f protoreflect.FieldDescriptor) bool {
		msg, ok := f.Parent().(protoreflect.MessageDescriptor)
		return ok && isRunRequestMessage(msg) && string(f.Name()) == "name"
	},
	LintField: func(f protoreflect.FieldDescriptor) []lint.Problem {
		// Rule check: Establish that the `resource_reference` annotation's
		// type ends in "Job".
		ref := utils.GetResourceReference(f)
		if ref != nil && !strings.HasSuffix(ref.GetType(), "Job") {
			suggestion := ""
			if msg, ok := f.Parent().(protoreflect.MessageDescriptor); ok {
				suggestion = strings.TrimPrefix(string(msg.Name()), "Run")
				suggestion = strings.TrimSuffix(suggestion, "Request")
			}

			return []lint.Problem{{
				Message:    fmt.Sprintf("The `type` of the `google.api.resource_reference` annotation should end in %q.", "Job"),
				Descriptor: f,
				Location:   locations.FieldResourceReference(f),
				Suggestion: suggestion,
			}}
		}
		return nil
	},
}
