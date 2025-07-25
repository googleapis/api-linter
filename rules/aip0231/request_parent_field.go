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

package aip0231

import (
	"fmt"
	"strings"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
	"github.com/stoewer/go-strcase"
)

// The Batch Get request message should have parent field.
var requestParentField = &lint.MessageRule{
	Name: lint.NewRuleName(231, "request-parent-field"),
	OnlyIf: func(m protoreflect.MessageDescriptor) bool {
		// In order to parse out the pattern, we get the resource message
		// from the response, then get the resource annotation from that,
		// and then inspect the pattern there (oy!).
		plural := strings.TrimPrefix(strings.TrimSuffix(string(m.Name()), "Request"), "BatchGet")
		if resp := utils.FindMessage(m.ParentFile(), fmt.Sprintf("BatchGet%sResponse", plural)); resp != nil {
			if resField := resp.Fields().ByName(protoreflect.Name(strcase.SnakeCase(plural))); resField != nil {
				if !utils.HasParent(utils.GetResource(resField.Message())) {
					return false
				}
			}
		}

		return isBatchGetRequestMessage(m)
	},
	LintMessage: utils.LintFieldPresentAndSingularString("parent"),
}
