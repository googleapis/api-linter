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

package aip0158

import (
	"fmt"

	"github.com/googleapis/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
)

var responseRepeatedFirstField = &lint.MessageRule{
	Name:   lint.NewRuleName(158, "response-repeated-first-field"),
	OnlyIf: isPaginatedResponseMessage,
	LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
		if len(m.GetFields()) == 0 {
			return []lint.Problem{{
				Message:    fmt.Sprintf("Paginated RPCs' response should have at least 1 field."),
			}}
		}

		if !m.FindFieldByNumber(1).IsRepeated() {
			return []lint.Problem{{
				Message:    fmt.Sprintf("First field by number of Paginated RPCs' response should be repeated."),
				Descriptor: m.FindFieldByNumber(1),
			}}
		}

		if !m.GetFields()[0].IsRepeated() {
			return []lint.Problem{{
				Message:    fmt.Sprintf("First field by position of Paginated RPCs' response should be repeated."),
				Descriptor: m.GetFields()[0],
			}}
		}

		return nil
	},
}
