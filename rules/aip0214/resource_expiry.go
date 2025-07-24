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

package aip0214

import (
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var resourceExpiry = &lint.FieldRule{
	Name: lint.NewRuleName(214, "resource-expiry"),
	OnlyIf: func(f protoreflect.FieldDescriptor) bool {
		isResource := utils.IsResource(f.Parent().(protoreflect.MessageDescriptor))
		isExpireTime := f.Name() == "expire_time"
		return isResource && isExpireTime
	},
	LintField: func(f protoreflect.FieldDescriptor) []lint.Problem {
		// If this field is output only, then there is no user input permitted
		// and therefore having a `ttl` field does not matter.
		if utils.GetFieldBehavior(f).Contains("OUTPUT_ONLY") {
			return nil
		}

		// If this message does not have a `ttl` field, suggest one.
		if f.GetOwner().FindFieldByName("ttl") == nil {
			return []lint.Problem{{
				Message:    "Resources that let users set expire_time should include an input only `ttl` field.",
				Descriptor: f,
			}}
		}

		return nil
	},
}
