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

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
	"github.com/stoewer/go-strcase"
)

// The resource field in a update method should named properly.
var requestResourceField = &lint.FieldRule{
	Name: lint.NewRuleName(134, "request-resource-field"),
	OnlyIf: func(f protoreflect.FieldDescriptor) bool {
		if message, ok := f.Parent().(protoreflect.MessageDescriptor); ok {
			return utils.IsUpdateRequestMessage(message) &&
				f.Message() != nil &&
				string(f.Message().Name()) == extractResource(string(message.Name()))
		}
		return false
	},
	LintField: func(f protoreflect.FieldDescriptor) []lint.Problem {
		resourceName := extractResource(string(f.Parent().(protoreflect.MessageDescriptor).Name()))
		wantFieldName := strcase.SnakeCase(resourceName)
		if string(f.Name()) != wantFieldName {
			return []lint.Problem{{
				Message:    fmt.Sprintf("Resource field should be named %q.", wantFieldName),
				Descriptor: f,
				Suggestion: wantFieldName,
				Location:   locations.DescriptorName(f),
			}}
		}
		return nil
	},
}