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
	"google.golang.org/protobuf/reflect/protoregistry"

	// These imports are for populating the global registry
	_ "cloud.google.com/go/longrunning/autogen/longrunningpb"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	_ "google.golang.org/genproto/googleapis/api/httpbody"
	_ "google.golang.org/genproto/googleapis/type/date"
	_ "google.golang.org/genproto/googleapis/type/datetime"
	_ "google.golang.org/genproto/googleapis/type/timeofday"
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

	// Create a custom resolver to ensure that we read from the global registry
	// which allows us to get the googleapis for our tests
	googleAPIRegistryResolver := protocompile.ResolverFunc(func(path string) (protocompile.SearchResult, error) {
		fd, err := protoregistry.GlobalFiles.FindFileByPath(path)
		if err != nil {
			return protocompile.SearchResult{}, err
		}
		return protocompile.SearchResult{Desc: fd}, nil
	})

	compiler := protocompile.Compiler{
		// Combine the in-memory resolver with the resolver for the global registry.
		// protocompile.WithStandardImports ensures that well-known types are also
		// resolved correctly.
		Resolver: protocompile.WithStandardImports(
			protocompile.CompositeResolver{
				memResolver,
				googleAPIRegistryResolver,
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
