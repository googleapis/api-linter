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

var resourcePattern = &lint.MessageRule{
	Name:   lint.NewRuleName(123, "resource-pattern"),
	OnlyIf: hasResourceAnnotation,
	LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
		resource := utils.GetResource(m)
		return lintResourcePattern(resource, m, locations.MessageResource(m))
	},
}

func lintResourcePattern(resource *annotations.ResourceDescriptor, desc desc.Descriptor, loc *dpb.SourceCodeInfo_Location) []lint.Problem {
	// Are any patterns declared at all? If not, complain.
	if len(resource.GetPattern()) == 0 {
		return []lint.Problem{{
			Message:    "Resources should declare resource name pattern(s).",
			Descriptor: desc,
			Location:   loc,
		}}
	}

	// Ensure that the constant segments of the pattern uses camel case,
	// not snake case, and there are no spaces.
	for _, pattern := range resource.GetPattern() {
		plainPattern := getPlainPattern(pattern)

		if strings.Contains(plainPattern, "_") {
			return []lint.Problem{{
				Message: fmt.Sprintf(
					"Resource patterns should use camel case (apart from the variable names), such as %q.",
					getDesiredPattern(pattern),
				),
				Descriptor: desc,
				Location:   loc,
			}}
		}
		if strings.Contains(plainPattern, " ") {
			return []lint.Problem{{
				Message:    "Resource patterns should not have spaces",
				Descriptor: desc,
				Location:   loc,
			}}
		}
	}
	return nil
}
