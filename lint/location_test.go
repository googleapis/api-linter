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

package lint

import (
	"strings"
	"testing"

	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/google/go-cmp/cmp"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/lithammer/dedent"
)

func TestLocations(t *testing.T) {
	parser := protoparse.Parser{
		Accessor: protoparse.FileContentsFromMap(map[string]string{
			"test.proto": strings.TrimSpace(dedent.Dedent(`
				// proto3 rules!
				syntax = "proto3";

				package google.api.linter;

				message Foo {
				  string bar = 1;
				}
			`)),
		}),
		IncludeSourceCodeInfo: true,
	}
	fds, err := parser.ParseFiles("test.proto")
	if err != nil {
		t.Fatalf("%v", err)
	}
	f := fds[0]

	// Test the file location functions.
	t.Run("File", func(t *testing.T) {
		tests := []struct {
			testName string
			fx       func(f *desc.FileDescriptor) *dpb.SourceCodeInfo_Location
			wantSpan []int32
		}{
			{"Syntax", SyntaxLocation, []int32{1, 0, int32(len("syntax = \"proto3\";"))}},
			{"Package", PackageLocation, []int32{3, 0, int32(len("package google.api.linter;"))}},
		}
		for _, test := range tests {
			t.Run(test.testName, func(t *testing.T) {
				if diff := cmp.Diff(test.fx(f).Span, test.wantSpan); diff != "" {
					t.Errorf(diff)
				}
			})
		}
	})

	// Test descriptor names.
	t.Run("DescriptorNames", func(t *testing.T) {
		tests := []struct {
			testName string
			d        desc.Descriptor
			wantSpan []int32
		}{
			{"Message", f.GetMessageTypes()[0], []int32{5, 8, 11}},
			{"Field", f.GetMessageTypes()[0].GetFields()[0], []int32{6, 9, 12}},
		}
		for _, test := range tests {
			t.Run(test.testName, func(t *testing.T) {
				if diff := cmp.Diff(DescriptorNameLocation(test.d).Span, test.wantSpan); diff != "" {
					t.Errorf(diff)
				}
			})
		}
	})

	// Test bogus locations.
	t.Run("Bogus", func(t *testing.T) {
		tests := []struct {
			testName string
			path     []int32
		}{
			{"NotFound", []int32{6, 0}},
		}
		for _, test := range tests {
			t.Run(test.testName, func(t *testing.T) {
				if loc := PathLocation(f, test.path); loc != nil {
					t.Errorf("%v", loc)
				}
			})
		}
	})
}

func TestMissingLocations(t *testing.T) {
	m, err := builder.NewMessage("Foo").Build()
	if err != nil {
		t.Fatalf("%v", err)
	}
	f := m.GetFile()
	tests := []struct {
		testName string
		fx       func(f *desc.FileDescriptor) *dpb.SourceCodeInfo_Location
	}{
		{"Syntax", SyntaxLocation},
		{"Package", PackageLocation},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			if diff := cmp.Diff(test.fx(f).Span, []int32{0, 0, 0}); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
