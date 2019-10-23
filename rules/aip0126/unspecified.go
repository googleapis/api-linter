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

package aip0126

import (
	"fmt"
	"strings"

	"github.com/googleapis/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
	"github.com/stoewer/go-strcase"
)

// This rule enforces that all enums have a default unspecified value,
// as described in [AIP-126](http://aip.dev/126).
//
// Because our APIs create automatically-generated client libraries, we need
// to consider languages that have varying behavior around default values.
// To avoid any ambiguity or confusion across languages, all enumerations
// should use an "unspecified" value beginning with the name of the enum
// itself as the first (0) value.
//
// ## Details
//
// This rule finds all enumerations and ensures that the first one is
// named after the enum itself with an `_UNSPECIFIED` suffix attached.
//
// ## Examples
//
// **Incorrect** code for this rule:
//
//   enum Format {
//     HARDCOVER = 0;  // Should have "FORMAT_UNSPECIFIED" first.
//   }
//
//   enum Format {
//     UNSPECIFIED = 0;  // Should be "FORMAT_UNSPECIFIED".
//     HARDCOVER = 1;
//   }
//
// **Correct** code for this rule:
//
//   enum Format {
//     FORMAT_UNSPECIFIED = 0;
//     HARDCOVER = 1;
//   }
//
// ## Disabling
//
// If you need to violate this rule, use a leading comment above the enum
// value.
//
//   enum Format {
//     // (-- api-linter: core::0126::unspecified=disabled --)
//     HARDCOVER = 0;
//   }
//
// If you need to violate this rule for an entire file, place the comment at
// the top of the file.
var unspecified = &lint.EnumRule{
	Name: lint.NewRuleName("core", "0126", "unspecified"),
	URI:  "https://aip.dev/126#guidance",
	LintEnum: func(e *desc.EnumDescriptor) []lint.Problem {
		firstValue := e.GetValues()[0]
		want := strings.ToUpper(strcase.SnakeCase(e.GetName()) + "_UNSPECIFIED")
		if firstValue.GetName() != want {
			return []lint.Problem{{
				Message:    fmt.Sprintf("The first enum value should be %q", want),
				Suggestion: want,
				Descriptor: firstValue,
				Location:   lint.DescriptorNameLocation(firstValue),
			}}
		}

		return nil
	},
}
