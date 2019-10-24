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
	"bytes"
	"fmt"
	"html/template"
	"strings"
	"testing"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/lithammer/dedent"

	// These imports cause the common protos to be registered with
	// the protocol buffer registry, and therefore make the call to
	// `proto.FileDescriptor` work for the imported files.
	_ "google.golang.org/genproto/googleapis/api/annotations"
	_ "google.golang.org/genproto/googleapis/longrunning"
)

// ParseProtoString parses a string representing a proto file, and returns
// a FileDescriptor.
//
// It dedents the string before parsing.
// It is unable to handle imports, and calls t.Fatalf if there is any error.
func ParseProtoString(t *testing.T, src string) *desc.FileDescriptor {
	// Make a map of filenames and file contents.
	// We hard-code "test.proto"; we do not care what the filename is.
	//
	// Include the common protos here too, so our proto may safely import them
	// if needed.
	fileContents := map[string]string{
		"test.proto": strings.TrimSpace(dedent.Dedent(src)),
	}

	// Parse the file.
	parser := protoparse.Parser{
		Accessor:              protoparse.FileContentsFromMap(fileContents),
		IncludeSourceCodeInfo: true,
		LookupImport:          desc.LoadFileDescriptor,
	}
	fds, err := parser.ParseFiles("test.proto")
	if err != nil {
		t.Fatalf("%v", err)
	}
	return fds[0]
}

// ParseProto3String parses a string representing a proto file, and returns
// a FileDescriptor.
//
// It adds the `syntax = "proto3";` line to the beginning of the file
// before calling ParseProtoString.
func ParseProto3String(t *testing.T, src string) *desc.FileDescriptor {
	return ParseProtoString(t, fmt.Sprintf(
		"syntax = \"proto3\";\n\n%s",
		strings.TrimSpace(dedent.Dedent(src)),
	))
}

// ParseProto3Tmpl parses a template string representing a proto file, and
// returns a FileDescriptor.
//
// It parses the template using Go's template Parse function, and then
// calls ParseProto3String.
func ParseProto3Tmpl(t *testing.T, src string, data interface{}) *desc.FileDescriptor {
	// Create a new template object.
	tmpl, err := template.New("test").Parse(src)
	if err != nil {
		t.Fatalf("Unable to parse Go template: %v", err)
	}

	// Execute the template and write the results to a bytes representing
	// the desired proto.
	var protoBytes bytes.Buffer
	err = tmpl.Execute(&protoBytes, data)
	if err != nil {
		t.Fatalf("Unable to execute Go template: %v", err)
	}

	// Parse the proto as a string.
	return ParseProto3String(t, protoBytes.String())
}
