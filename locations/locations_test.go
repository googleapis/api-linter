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
	"strings"
	"testing"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/lithammer/dedent"

	// These imports cause the common protos to be registered with
	// the protocol buffer registry, and therefore make the call to
	// `proto.FileDescriptor` work for the imported files.
	_ "cloud.google.com/go/longrunning/autogen/longrunningpb"
	_ "google.golang.org/genproto/googleapis/api/annotations"
)

func parse(t *testing.T, s string) *desc.FileDescriptor {
	s = strings.TrimSpace(dedent.Dedent(s))
	if !strings.Contains(s, "syntax = ") {
		s = "syntax = \"proto3\";\n\n" + s
	}
	parser := protoparse.Parser{
		Accessor: protoparse.FileContentsFromMap(map[string]string{
			"test.proto": strings.TrimSpace(dedent.Dedent(s)),
		}),
		IncludeSourceCodeInfo: true,
		LookupImport:          desc.LoadFileDescriptor,
	}
	fds, err := parser.ParseFiles("test.proto")
	if err != nil {
		t.Fatalf("%v", err)
	}
	return fds[0]
}
