// Copyright 2023 Google LLC
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

package aip0148

import (
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/types/descriptorpb"
)

var uidFormat = &lint.FieldRule{
	Name: lint.NewRuleName(148, "uid-format"),
	OnlyIf: func(fd *desc.FieldDescriptor) bool {
		return fd.GetType() == descriptorpb.FieldDescriptorProto_TYPE_STRING && fd.GetName() == uidStr
	},
	LintField: func(fd *desc.FieldDescriptor) []lint.Problem {
		if !utils.HasFormat(fd) || utils.GetFormat(fd) != annotations.FieldInfo_UUID4 {
			return []lint.Problem{{
				Message:    "The `uid` field must have a `(google.api.field_info).format = UUID4` annotation.",
				Descriptor: fd,
			}}
		}

		return nil
	},
}
