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

package aip0133

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
	"github.com/stoewer/go-strcase"
)

var methodSignature = &lint.MethodRule{
	Name:   lint.NewRuleName(133, "method-signature"),
	OnlyIf: isCreateMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		signatures := utils.GetMethodSignatures(m)

		// Determine what signature we want. The {resource}_id is desired
		// if and only if the field exists on the request.
		resourceField := strcase.SnakeCase(getResourceMsgName(m))
		want := []string{}
		if !hasNoParent(m.GetOutputType()) {
			want = append(want, "parent")
		}
		want = append(want, resourceField)
		if idField := resourceField + "_id"; m.GetInputType().FindFieldByName(idField) != nil {
			want = append(want, idField)
		}

		// Check if the signature is missing.
		if len(signatures) == 0 {
			return []lint.Problem{{
				Message: fmt.Sprintf(
					"Create methods should include `(google.api.method_signature) = %q`",
					strings.Join(want, ","),
				),
				Descriptor: m,
			}}
		}

		// Check if the signature is wrong.
		if !reflect.DeepEqual(signatures[0], want) {
			return []lint.Problem{{
				Message: fmt.Sprintf(
					"The method signature for Create methods should be %q.",
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
