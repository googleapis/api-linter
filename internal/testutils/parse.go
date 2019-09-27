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

package testutils

import (
	"strings"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/lithammer/dedent"
)

// ParseProtoString parses a string representing a proto file, and returns
// a FileDescriptor.
//
// It dedents the string before parsing.
// It is unable to handle imports, and panics if there is any error.
func ParseProtoString(src string) *desc.FileDescriptor {
	parser := protoparse.Parser{
		Accessor: protoparse.FileContentsFromMap(map[string]string{
			"test.proto": strings.TrimSpace(dedent.Dedent(src)),
		}),
		IncludeSourceCodeInfo: true,
	}
	fds, err := parser.ParseFiles("test.proto")
	if err != nil {
		panic(err)
	}
	return fds[0]
}
