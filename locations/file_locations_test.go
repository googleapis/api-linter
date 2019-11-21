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
	"testing"

	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/google/go-cmp/cmp"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"
)

func TestLocations(t *testing.T) {
	f := parse(t, `
		// proto3 rules!
		syntax = "proto3";

		package google.api.linter;

		message Foo {
			string bar = 1;
		}
	`)

	// Test the file location functions.
	t.Run("File", func(t *testing.T) {
		tests := []struct {
			testName string
			fx       func(f *desc.FileDescriptor) *dpb.SourceCodeInfo_Location
			wantSpan []int32
		}{
			{"Syntax", FileSyntax, []int32{1, 0, int32(len("syntax = \"proto3\";"))}},
			{"Package", FilePackage, []int32{3, 0, int32(len("package google.api.linter;"))}},
		}
		for _, test := range tests {
			t.Run(test.testName, func(t *testing.T) {
				if diff := cmp.Diff(test.fx(f).Span, test.wantSpan); diff != "" {
					t.Errorf(diff)
				}
			})
		}
	})

	// Test bogus locations.
	t.Run("Bogus", func(t *testing.T) {
		tests := []struct {
			testName string
			path     []int
		}{
			{"NotFound", []int{6, 0}},
		}
		for _, test := range tests {
			t.Run(test.testName, func(t *testing.T) {
				if loc := pathLocation(f, test.path...); loc != nil {
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
		{"Syntax", FileSyntax},
		{"Package", FilePackage},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			if diff := cmp.Diff(test.fx(f).Span, []int32{0, 0, 0}); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
