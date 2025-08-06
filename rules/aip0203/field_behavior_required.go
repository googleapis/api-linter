// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aip0203

import (
	"fmt"

	"bitbucket.org/creachadair/stringset"
	"github.com/googleapis/api-linter/v2/lint"
	"github.com/googleapis/api-linter/v2/rules/internal/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var minimumRequiredFieldBehavior = stringset.New(
	"OPTIONAL", "REQUIRED", "OUTPUT_ONLY", "IMMUTABLE",
)

var excusedResourceFields = stringset.New(
	"name", // Uses https://google.aip.dev/203#identifier
	"etag", // Prohibited by https://google.aip.dev/154
)

var fieldBehaviorRequired = &lint.MethodRule{
	Name: lint.NewRuleName(203, "field-behavior-required"),
	LintMethod: func(m protoreflect.MethodDescriptor) []lint.Problem {
		req := m.Input()
		p := string(m.ParentFile().Package())
		ps := problems(req, p, map[protoreflect.Descriptor]bool{})
		if len(ps) == 0 {
			return nil
		}

		return ps
	},
}

func problems(m protoreflect.MessageDescriptor, pkg string, visited map[protoreflect.Descriptor]bool) []lint.Problem {
	var ps []lint.Problem

	// Ensure the input type, or recursively visited message, is part of the
	// same package before linting.
	if string(m.ParentFile().Package()) != pkg {
		return nil
	}

	for i := 0; i < m.Fields().Len(); i++ {
		f := m.Fields().Get(i)
		// ignore the field if it was already visited
		if ok := visited[f]; ok {
			continue
		}
		visited[f] = true

		if utils.IsResource(m) && excusedResourceFields.Contains(string(f.Name())) {
			continue
		}

		// Ignore a field if it is a OneOf (do not ignore children)
		if f.ContainingOneof() == nil {
			p := checkFieldBehavior(f)
			if p != nil {
				ps = append(ps, *p)
			}
		}

		if mt := f.Message(); mt != nil && !mt.IsMapEntry() {
			ps = append(ps, problems(mt, pkg, visited)...)
		}
	}

	return ps
}

func checkFieldBehavior(f protoreflect.FieldDescriptor) *lint.Problem {
	fb := utils.GetFieldBehavior(f)

	if len(fb) == 0 {
		return &lint.Problem{
			Message:    fmt.Sprintf("google.api.field_behavior annotation must be set on %q and contain one of, \"%v\"", f.Name(), minimumRequiredFieldBehavior),
			Descriptor: f,
		}
	}

	if !minimumRequiredFieldBehavior.Intersects(fb) {
		// check for at least one valid annotation
		return &lint.Problem{
			Message: fmt.Sprintf(
				"google.api.field_behavior on field %q must contain at least one, \"%v\"", f.Name(), minimumRequiredFieldBehavior),
			Descriptor: f,
		}
	}

	return nil
}
