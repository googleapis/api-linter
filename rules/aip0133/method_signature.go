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

	"github.com/googleapis/api-linter/v2/lint"
	"github.com/googleapis/api-linter/v2/locations"
	"github.com/googleapis/api-linter/v2/rules/internal/utils"
	"github.com/stoewer/go-strcase"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var methodSignature = &lint.MethodRule{
	Name: lint.NewRuleName(133, "method-signature"),
	OnlyIf: func(m protoreflect.MethodDescriptor) bool {
		return utils.IsCreateMethod(m) && utils.IsResource(utils.GetResponseType(m))
	},
	LintMethod: func(m protoreflect.MethodDescriptor) []lint.Problem {
		signatures := utils.GetMethodSignatures(m)

		// Determine what signature we want.
		want := []string{}
		if utils.HasParent(utils.GetResource(utils.GetResponseType(m))) {
			want = append(want, "parent")
		}
		for i := 0; i < m.Input().Fields().Len(); i++ {
			f := m.Input().Fields().Get(i)
			if mt := f.Message(); mt != nil && utils.IsResource(mt) {
				want = append(want, string(f.Name()))
				break
			}
		}
		// The {resource}_id is desired if and only if the field exists on the
		// request and the request targets a resource.
		expectedResourceIDField := strcase.SnakeCase(utils.GetResourceMessageName(m, "Create"))
		if idField := expectedResourceIDField + "_id"; m.Input().Fields().ByName(protoreflect.Name(idField)) != nil {
			want = append(want, idField)
		}

		// The Standard Create is not standard and has nothing to suggest.
		// There are likely other rules warning about the non-standard nature
		// so just silently move on.
		if len(want) == 0 {
			return nil
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
