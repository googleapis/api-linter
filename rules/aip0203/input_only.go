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

package aip0203

import (
	"regexp"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/genproto/googleapis/api/annotations"
)

var inputOnly = &lint.FieldRule{
	Name:   lint.NewRuleName(203, "input-only"),
	OnlyIf: withoutInputOnlyFieldBehavior,
	LintField: func(f *desc.FieldDescriptor) []lint.Problem {
		return checkLeadingComments(f, inputOnlyRegexp, "INPUT_ONLY")
	},
}

var inputOnlyRegexp = regexp.MustCompile("(?i).*input.?only.*")

func withoutInputOnlyFieldBehavior(f *desc.FieldDescriptor) bool {
	for _, v := range utils.GetFieldBehavior(f) {
		if v == annotations.FieldBehavior_INPUT_ONLY {
			return false
		}
	}
	return true
}
