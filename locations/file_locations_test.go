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

	"github.com/google/go-cmp/cmp"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"
	dpb "google.golang.org/protobuf/types/descriptorpb"
)

func TestLocations(t *testing.T) {
	f := parse(t, `
		// proto3 rules!
		syntax = "proto3";

		import "google/api/resource.proto";

		package google.api.linter;

		option csharp_namespace = "Google.Api.Linter";
		option java_package = "com.google.api.linter";
		option php_namespace = "Google\\Api\\Linter";
		option ruby_package = "Google::Api::Linter";
		option cc_enable_arenas = false;

		message Foo {
			string bar = 1;
		}
	`)

	// Test the file location functions.
	t.Run("File", func(t *testing.T) {
		tests := []struct {
			testName string
			fx       func(f *desc.FileDescriptor) *dpb.SourceCodeInfo_Location
			idxFx    func(f *desc.FileDescriptor, i int) *dpb.SourceCodeInfo_Location
			idx      int
			wantSpan []int32
		}{
			{
				testName: "Syntax",
				fx:       FileSyntax,
				wantSpan: []int32{1, 0, int32(len("syntax = \"proto3\";"))},
			},
			{
				testName: "Package",
				fx:       FilePackage,
				wantSpan: []int32{5, 0, int32(len("package google.api.linter;"))},
			},
			{
				testName: "CsharpNamespace",
				fx:       FileCsharpNamespace,
				wantSpan: []int32{7, 0, int32(len(`option csharp_namespace = "Google.Api.Linter";`))},
			},
			{
				testName: "JavaPackage",
				fx:       FileJavaPackage,
				wantSpan: []int32{8, 0, int32(len(`option java_package = "com.google.api.linter";`))},
			},
			{
				testName: "PhpNamespace",
				fx:       FilePhpNamespace,
				wantSpan: []int32{9, 0, int32(len(`option php_namespace = "Google\\Api\\Linter";`))},
			},
			{
				testName: "RubyPackage",
				fx:       FileRubyPackage,
				wantSpan: []int32{10, 0, int32(len(`option ruby_package = "Google::Api::Linter";`))},
			},
			{
				testName: "Import",
				idxFx:    FileImport,
				idx:      0,
				wantSpan: []int32{3, 0, int32(len(`import "google/api/resource.proto";`))},
			},
			{
				testName: "CCEnableArenas",
				fx:       FileCCEnableArenas,
				wantSpan: []int32{11, 0, int32(len(`option cc_enable_arenas = false;`))},
			},
		}
		for _, test := range tests {
			t.Run(test.testName, func(t *testing.T) {
				var l *dpb.SourceCodeInfo_Location
				if test.fx != nil {
					l = test.fx(f)
				} else {
					l = test.idxFx(f, test.idx)
				}
				if diff := cmp.Diff(l.Span, test.wantSpan); diff != "" {
					t.Error(diff)
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
				t.Error(diff)
			}
		})
	}
}
