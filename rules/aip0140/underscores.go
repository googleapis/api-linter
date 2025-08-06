// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aip0140

import (
	"regexp"
	"strings"

	"github.com/googleapis/api-linter/v2/lint"
	"github.com/googleapis/api-linter/v2/locations"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var underscores = &lint.FieldRule{
	Name: lint.NewRuleName(140, "underscores"),
	LintField: func(f protoreflect.FieldDescriptor) []lint.Problem {
		n := string(f.Name())
		if strings.HasPrefix(n, "_") || strings.HasSuffix(n, "_") || strings.Contains(n, "__") {
			return []lint.Problem{{
				Message:    "Field names must not begin or end with underscore, or have adjacent underscores.",
				Descriptor: f,
				Location:   locations.DescriptorName(f),
				Suggestion: adjacent.ReplaceAllString(strings.TrimRight(strings.TrimLeft(n, "_"), "_"), "_"),
			}}
		}
		return nil
	},
}

var adjacent = regexp.MustCompile("_{2,}")
