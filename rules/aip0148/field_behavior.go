// Copyright 2021 Google LLC
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
	"bitbucket.org/creachadair/stringset"
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var fieldBehavior = &lint.FieldRule{
	Name: lint.NewRuleName(148, "field-behavior"),
	OnlyIf: func(f protoreflect.FieldDescriptor) bool {
		return utils.IsResource(f.GetOwner()) && outputOnlyFields.Contains(f.Name())
	},
	LintField: utils.LintOutputOnlyField,
}

var outputOnlyFields = stringset.New(
	"create_time",
	"delete_time",
	"uid",
	"update_time",
)
