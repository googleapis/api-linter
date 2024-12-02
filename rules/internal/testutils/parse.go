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
	"strings"
	"testing"
	"text/template"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/lithammer/dedent"

	// These imports cause the common protos to be registered with
	// the protocol buffer registry, and therefore make the call to
	// `proto.FileDescriptor` work for the imported files.
	_ "cloud.google.com/go/longrunning/autogen/longrunningpb"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	_ "google.golang.org/genproto/googleapis/api/httpbody"
	_ "google.golang.org/genproto/googleapis/type/date"
	_ "google.golang.org/genproto/googleapis/type/datetime"
	_ "google.golang.org/genproto/googleapis/type/timeofday"
)

// ParseProtoStrings parses a map representing a proto files, and returns
// a slice of FileDescriptors.
//
// It dedents the string before parsing.
func ParseProtoStrings(t *testing.T, src map[string]string) map[string]*desc.FileDescriptor {
	filenames := []string{}
	for k, v := range src {
		filenames = append(filenames, k)
		src[k] = strings.TrimSpace(dedent.Dedent(v))
	}

	// Parse the file.
	parser := protoparse.Parser{
		Accessor:              protoparse.FileContentsFromMap(src),
		IncludeSourceCodeInfo: true,
		LookupImport:          desc.LoadFileDescriptor,
	}
	fds, err := parser.ParseFiles(filenames...)
	if err != nil {
		t.Fatalf("%v", err)
	}
	answer := map[string]*desc.FileDescriptor{}
	for _, fd := range fds {
		answer[fd.GetName()] = fd
	}
	return answer
}

// ParseProto3String parses a string representing a proto file, and returns
// a FileDescriptor.
//
// It adds the `syntax = "proto3";` line to the beginning of the file and
// chooses a filename, and then calls ParseProtoStrings.
func ParseProto3String(t *testing.T, src string) *desc.FileDescriptor {
	return ParseProtoStrings(t, map[string]string{
		"test.proto": fmt.Sprintf(
			"syntax = \"proto3\";\n\n%s",
			strings.TrimSpace(dedent.Dedent(src)),
		),
	})["test.proto"]
}

// ParseProto3Tmpl parses a template string representing a proto file, and
// returns a FileDescriptor.
//
// It parses the template using Go's text/template Parse function, and then
// calls ParseProto3String.
func ParseProto3Tmpl(t *testing.T, src string, data interface{}) *desc.FileDescriptor {
	return ParseProto3Tmpls(t, map[string]string{
		"test.proto": src,
	}, data)["test.proto"]
}

// ParseProto3Tmpls parses template strings representing a proto file,
// and returns FileDescriptors.
//
// It parses the template using Go's text/template Parse function, and then
// calls ParseProto3Strings.
func ParseProto3Tmpls(t *testing.T, srcs map[string]string, data interface{}) map[string]*desc.FileDescriptor {
	strs := map[string]string{}
	for fn, src := range srcs {
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

		// Add the proto to the map to send to parse strings.
		strs[fn] = fmt.Sprintf("syntax = %q;\n\n%s", "proto3", protoBytes.String())
	}

	// Parse the proto as a string.
	return ParseProtoStrings(t, strs)
}
