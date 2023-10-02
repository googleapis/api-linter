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

package aip0202

import (
	"fmt"

	"bitbucket.org/creachadair/stringset"
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
	apb "google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/types/descriptorpb"
)

var stringOnlyFormats = stringset.New(
	apb.FieldInfo_IPV4.String(),
	apb.FieldInfo_IPV6.String(),
	apb.FieldInfo_IPV4_OR_IPV6.String(),
	apb.FieldInfo_UUID4.String(),
)

var stringOnlyFormat = &lint.FieldRule{
	Name: lint.NewRuleName(202, "string-only-format"),
	OnlyIf: func(fd *desc.FieldDescriptor) bool {
		return utils.HasFormat(fd) && fd.GetType() != descriptorpb.FieldDescriptorProto_TYPE_STRING
	},
	LintField: func(fd *desc.FieldDescriptor) []lint.Problem {
		// Field being linted is not a string, check that it isn't using a string-only format.
		if format := utils.GetFormat(fd).String(); stringOnlyFormats.Contains(format) {
			return []lint.Problem{{
				Message:    fmt.Sprintf("Format %q must only be used on string fields", format),
				Descriptor: fd,
			}}
		}
		return nil
	},
}
