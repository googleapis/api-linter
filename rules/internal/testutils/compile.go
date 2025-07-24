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

package testutils

import (
	"context"
	"testing"

	"github.com/bufbuild/protocompile"
	"github.com/bufbuild/protocompile/linker"
)

var background = context.Background()

// Compile compiles a single proto file from a string.
func Compile(t *testing.T, template string, data any) linker.File {
	t.Helper()
	content := ParseTemplate(t, template, data)
	memResolver := &protocompile.SourceResolver{
		Accessor: protocompile.SourceAccessorFromMap(map[string]string{
			"test.proto": "syntax = \"proto3\";\n" + content,
		}),
	}

	compiler := protocompile.Compiler{
		Resolver: protocompile.WithStandardImports(memResolver),
	}
	files, err := compiler.Compile(
		background,
		"test.proto",
	)
	if err != nil {
		t.Fatalf("Failed to compile: %v", err)
	}
	return files[0]
}

// CompileStrings compiles multiple proto files from a map of strings.
func CompileStrings(t *testing.T, files map[string]string) linker.Files {
	t.Helper()
	var fileNames []string
	for name := range files {
		fileNames = append(fileNames, name)
	}
	memResolver := &protocompile.SourceResolver{
		Accessor: protocompile.SourceAccessorFromMap(files),
	}
	compiler := protocompile.Compiler{
		Resolver: protocompile.WithStandardImports(memResolver),
	}
	compiled, err := compiler.Compile(
		context.Background(),
		fileNames...,
	)
	if err != nil {
		t.Fatalf("Failed to compile: %v", err)
	}
	return compiled
}
