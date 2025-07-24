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
	// TODO: this is a temp measurement. this requires linking the googleapis locally
	// and  it also requires depending on them. We might want to figure out a better
	// approach. maybe we could add it as a git submodule
	// Create a resolver for googleapis (used for standard API imports).
	// This assumes the `googleapis` directory exists relative to the project root.
	googleApisResolver := &protocompile.SourceResolver{
		ImportPaths: []string{"googleapis"},
	}

	compiler := protocompile.Compiler{
		// Combine the in-memory resolver with the googleapis resolver.
		// protocompile.WithStandardImports ensures that well-known types are also
		// resolved correctly.
		Resolver: protocompile.WithStandardImports(
			protocompile.CompositeResolver{
				memResolver,
				googleApisResolver,
			},
		),
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
