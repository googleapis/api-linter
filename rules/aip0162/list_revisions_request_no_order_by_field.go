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

package aip0162

import (
	"github.com/googleapis/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
)

// List Revisions requests should not have an order_by field.
var listRevisionsRequestNoOrderByField = &lint.FieldRule{
	Name: lint.NewRuleName(162, "list-revisions-request-no-order-by-field"),
	OnlyIf: func(f *desc.FieldDescriptor) bool {
		return IsListRevisionsRequestMessage(f.GetOwner()) && f.GetName() == "order_by"
	},
	LintField: func(f *desc.FieldDescriptor) []lint.Problem {
		return []lint.Problem{{
			Message:    "List Revisions requests should not contain an `order_by` field, as revisions must be ordered in reverse chronological order.",
			Descriptor: f,
		}}
	},
}
