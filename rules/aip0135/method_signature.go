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

package aip0135

import (
	"fmt"
	"strings"

	"bitbucket.org/creachadair/stringset"
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

var methodSignature = &lint.MethodRule{
	Name:   lint.NewRuleName(135, "method-signature"),
	OnlyIf: utils.IsDeleteMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		signatures := utils.GetMethodSignatures(m)
		in := m.GetInputType()

		fields := []string{"name"}
		if etag := in.FindFieldByName("etag"); etag != nil {
			fields = append(fields, etag.GetName())
		}
		if force := in.FindFieldByName("force"); force != nil {
			fields = append(fields, force.GetName())
		}
		want := strings.Join(fields, ",")

		// Check if the signature is missing.
		if len(signatures) == 0 {
			return []lint.Problem{{
				Message: fmt.Sprintf(
					"Delete methods should include `(google.api.method_signature) = %q`",
					want,
				),
				Descriptor: m,
			}}
		}

		// Check if the signature contains a disallowed field or doesn't contain
		// "name".
		first := signatures[0]
		fieldSet := stringset.New(fields...)
		if !fieldSet.Contains(first...) || !stringset.New(first...).Contains("name") {
			return []lint.Problem{{
				Message:    fmt.Sprintf("The method signature for Delete methods should be %q.", want),
				Suggestion: fmt.Sprintf("option (google.api.method_signature) = %q;", want),
				Descriptor: m,
				Location:   locations.MethodSignature(m, 0),
			}}
		}

		return nil
	},
}
