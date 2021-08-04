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

package aip0191

import (
	"path/filepath"
	"strings"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/jhump/protoreflect/desc"
)

// Protobuf package must match the directory structure.
var protoPkg = &lint.FileRule{
	Name: lint.NewRuleName(191, "proto-package"),
	LintFile: func(f *desc.FileDescriptor) []lint.Problem {
		dir := filepath.Dir(f.GetName())
		pkg := strings.ReplaceAll(f.GetPackage(), ".", string(filepath.Separator))

		if dir != "." && dir != pkg {
			return []lint.Problem{{
				Message:    "Proto package and directory structure mismatch: The proto package must match the proto directory structure.",
				Descriptor: f,
				Location:   locations.FilePackage(f),
			}}
		}

		return nil
	},
}
