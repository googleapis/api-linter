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

package aip0123

import (
	"fmt"
	"strings"

	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/genproto/googleapis/api/annotations"
)

var resourceVariables = &lint.MessageRule{
	Name:   lint.NewRuleName(123, "resource-variables"),
	OnlyIf: hasResourceAnnotation,
	LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
		resource := utils.GetResource(m)

		return lintResourceVariables(resource, m, locations.MessageResource(m))
	},
}

// lintResourceVariables lints the resource ID segments of the pattern(s) in the
// give ResourceDescriptor. This is used for both the file-level annotation
// google.api.resource_definition and the message-level annotation
// google.api.resource.
func lintResourceVariables(resource *annotations.ResourceDescriptor, desc desc.Descriptor, loc *dpb.SourceCodeInfo_Location) []lint.Problem {
	for _, pattern := range resource.GetPattern() {
		for _, variable := range getVariables(pattern) {
			if strings.ToLower(variable) != variable {
				return []lint.Problem{{
					Message: fmt.Sprintf(
						"Variable names in patterns should use snake case, such as %q.",
						getDesiredPattern(pattern),
					),
					Descriptor: desc,
					Location:   loc,
				}}
			}
			if strings.HasSuffix(variable, "_id") {
				return []lint.Problem{{
					Message: fmt.Sprintf(
						"Variable names should omit the `_id` suffix, such as %q.",
						getDesiredPattern(pattern),
					),
					Descriptor: desc,
					Location:   loc,
				}}
			}
		}
	}
	return nil
}
