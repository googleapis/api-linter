// Copyright 2019 Google LLC
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

package locations

import (
	"context"
	"strings"
	"sync"
	"testing"

	"github.com/bufbuild/protocompile"
	"github.com/lithammer/dedent"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"

	// These imports cause the common protos to be registered with
	// the protocol buffer registry, and therefore make the call to
	// `proto.FileDescriptor` work for the imported files.
	_ "cloud.google.com/go/longrunning/autogen/longrunningpb"
	_ "google.golang.org/genproto/googleapis/api/annotations"
)

func parse(t *testing.T, s string) protoreflect.FileDescriptor {
	t.Helper()
	s = strings.TrimSpace(dedent.Dedent(s))
	if !strings.Contains(s, "syntax = ") {
		s = "syntax = \"proto3\";\n\n" + s
	}

	// Resolver for our in-memory test file
	testFileResolver := &protocompile.SourceResolver{
		Accessor: protocompile.SourceAccessorFromMap(map[string]string{
			"test.proto": s,
		}),
	}

	// Resolver for standard imports (like google/api/annotations.proto)
	importResolver := protocompile.ResolverFunc(func(path string) (protocompile.SearchResult, error) {
		fd, err := protoregistry.GlobalFiles.FindFileByPath(path)
		if err != nil {
			return protocompile.SearchResult{}, err
		}
		return protocompile.SearchResult{Desc: fd}, nil
	})

	compiler := protocompile.Compiler{
		Resolver:       protocompile.CompositeResolver{testFileResolver, importResolver},
		SourceInfoMode: protocompile.SourceInfoStandard,
	}

	fds, err := compiler.Compile(context.Background(), "test.proto")
	if err != nil {
		t.Fatalf("%v", err)
	}

	return fds[0]
}

func TestSourceInfo_Concurrency(t *testing.T) {
	fd := parse(t, `
	syntax = "proto3";
	package foo.bar;
	`)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		FileSyntax(fd)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		FilePackage(fd)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		FileImport(fd, 0)
	}()
	wg.Wait()
}
