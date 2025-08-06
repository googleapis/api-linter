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

package aip0235

import (
	"fmt"

	"github.com/gertd/go-pluralize"
	"github.com/googleapis/api-linter/v2/lint"
	"github.com/googleapis/api-linter/v2/locations"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// The Batch Delete standard method should have repeated name field or repeated
// standard delete request message field, but the latter one is not suggested.
var requestNamesField = &lint.MessageRule{
	Name:   lint.NewRuleName(235, "request-names-field"),
	OnlyIf: isBatchDeleteRequestMessage,
	LintMessage: func(m protoreflect.MessageDescriptor) (problems []lint.Problem) {
		// Rule check: Establish that a name field is present.
		names := m.Fields().ByName("names")
		deleteReqMsg := m.Fields().ByName("requests")

		// Rule check: Ensure that the names field is present.
		if names == nil && deleteReqMsg == nil {
			problems = append(problems, lint.Problem{
				Message:    fmt.Sprintf(`Message %q has no "names" field`, m.Name()),
				Descriptor: m,
			})
		}

		// Rule check: Ensure that only the suggested names field is present.
		if names != nil && deleteReqMsg != nil {
			problems = append(problems, lint.Problem{
				Message:    fmt.Sprintf(`Message %q should delete "requests" field, only keep the "names" field`, m.Name()),
				Descriptor: deleteReqMsg,
			})
		}

		// Rule check: Ensure that the names field is repeated.
		if names != nil && !names.IsList() {
			problems = append(problems, lint.Problem{
				Message:    `The "names" field should be repeated`,
				Suggestion: "repeated string",
				Descriptor: names,
				Location:   locations.FieldType(names),
			})
		}

		// Rule check: Ensure that the names field is the correct type.
		if names != nil && names.Kind() != protoreflect.StringKind {
			problems = append(problems, lint.Problem{
				Message:    `"names" field on Batch Delete Request should be a "string" type`,
				Suggestion: "string",
				Descriptor: names,
				Location:   locations.FieldType(names),
			})
		}

		// Rule check: Ensure that the standard delete request message field is repeated.
		if deleteReqMsg != nil && !deleteReqMsg.IsList() {
			problems = append(problems, lint.Problem{
				Message:    `The "requests" field should be repeated`,
				Descriptor: deleteReqMsg,
			})
		}

		// Rule check: Ensure that the standard delete request message field is the
		// correct type. Note: Use m.Name()[11:len(m.Name())-7]) to retrieve
		// the resource name from the batch delete request, for example:
		// "BatchDeleteBooksRequest" -> "Books"
		rightTypeName := fmt.Sprintf("Delete%sRequest", pluralize.NewClient().Singular(string(m.Name())[11:len(m.Name())-7]))
		if deleteReqMsg != nil && (deleteReqMsg.Message() == nil || deleteReqMsg.Message().Name() != protoreflect.Name(rightTypeName)) {
			problems = append(problems, lint.Problem{
				Message:    fmt.Sprintf(`The "requests" field on Batch Delete Request should be a %q type`, rightTypeName),
				Descriptor: deleteReqMsg,
				Location:   locations.FieldType(deleteReqMsg),
			})
		}
		return
	},
}
