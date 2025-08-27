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

package aip0134

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/googleapis/api-linter/v2/lint"
	"github.com/googleapis/api-linter/v2/locations"
	"github.com/googleapis/api-linter/v2/rules/internal/utils"
	"github.com/stoewer/go-strcase"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var methodSignature = &lint.MethodRule{
	Name:   lint.NewRuleName(134, "method-signature"),
	OnlyIf: utils.IsUpdateMethod,
	LintMethod: func(m protoreflect.MethodDescriptor) []lint.Problem {
		signatures := utils.GetMethodSignatures(m)
		want := []string{
			strcase.SnakeCase(strings.TrimPrefix(string(m.Name()), "Update")),
			"update_mask",
		}

		// Check if the signature is missing.
		if len(signatures) == 0 {
			return []lint.Problem{{
				Message: fmt.Sprintf(
					"Update methods should include `(google.api.method_signature) = %q`",
					strings.Join(want, ","),
				),
				Descriptor: m,
			}}
		}

		// Check if the signature is wrong.
		if !reflect.DeepEqual(signatures[0], want) {
			return []lint.Problem{{
				Message: fmt.Sprintf(
					"The method signature for Update methods should be %q.",
					strings.Join(want, ","),
				),
				Suggestion: fmt.Sprintf(
					"option (google.api.method_signature) = %q;",
					strings.Join(want, ","),
				),
				Descriptor: m,
				Location:   locations.MethodSignature(m, 0),
			}}
		}
		return nil
	},
}
