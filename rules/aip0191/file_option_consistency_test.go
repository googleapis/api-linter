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

package aip0191

import (
	"reflect"
	"sort"
	"testing"

	"github.com/commure/api-linter/lint"
	"github.com/commure/api-linter/rules/internal/testutils"
	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"
	"github.com/stoewer/go-strcase"
	"google.golang.org/protobuf/proto"
)

func TestFileOptionConsistency(t *testing.T) {
	// Make a "control file".
	// This is a situation where using a builder is much easier
	// than parsing a proto template.
	controlFile := builder.NewFile("control.proto").SetPackageName(
		"google.example.v1",
	).SetOptions(getOptions(nil))

	// Run our tests. Each test will make a "test file" that imports the
	// control file.
	for _, test := range []consistencyTest{
		{"ValidSame", map[string]string{}},
		{"InvalidCsharpNamespaceMissing", map[string]string{"csharp_namespace": ""}},
		{"InvalidCsharpNamespaceMismatch", map[string]string{"csharp_namespace": "Example.V1"}},
		{"InvalidJavaPackageMissing", map[string]string{"java_package": ""}},
		{"InvalidJavaPackageMismatch", map[string]string{"java_package": "com.example.v1"}},
		{"InvalidPhpNamespaceMissing", map[string]string{"php_namespace": ""}},
		{"InvalidPhpNamespaceMismatch", map[string]string{"php_namespace": "Example\\V1"}},
		{"InvalidPhpMetadataNamespace", map[string]string{"php_metadata_namespace": "Example\\V1"}},
		{"InvalidPhpClassPrefix", map[string]string{"php_class_prefix": "ExampleProto"}},
		{"InvalidObjcClassPrefixMissing", map[string]string{"objc_class_prefix": ""}},
		{"InvalidObjcClassPrefixMismatch", map[string]string{"objc_class_prefix": "GEXV1"}},
		{"InvalidRubyPackageMissing", map[string]string{"ruby_package": ""}},
		{"InvalidRubyPackageMismatch", map[string]string{"ruby_package": "Example::V1"}},
		{"InvalidSwiftPrefixMissing", map[string]string{"swift_prefix": ""}},
		{"InvalidSwiftPrefixMismatch", map[string]string{"swift_prefix": "ExampleProto"}},
		{"InvalidLotsOfReasons", map[string]string{
			"csharp_namespace": "Example.V1",
			"go_package":       "google.golang.org/googleapis/genproto/googleapis/example/v1;example",
			"java_package":     "com.example.v1",
			"ruby_package":     "Example::V1",
			"swift_prefix":     "",
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			// Build our test file, which imports the control file.
			testFile, err := builder.NewFile("test.proto").AddDependency(
				controlFile,
			).SetPackageName("google.example.v1").SetOptions(getOptions(test.options)).Build()
			if err != nil {
				t.Fatalf("Could not build test file.")
			}
			if diff := test.getProblems(testFile).Diff(fileOptionConsistency.Lint(testFile)); diff != "" {
				t.Errorf(diff)
			}
		})
	}

	// Also ensure separate packages are ignored.
	t.Run("ValidDifferentPackages", func(t *testing.T) {
		testFile, err := builder.NewFile("test.proto").AddDependency(controlFile).SetOptions(
			getOptions(map[string]string{"java_package": "com.wrong"}),
		).SetPackageName("google.different.v1").Build()
		if err != nil {
			t.Fatalf("Could not build test file.")
		}
		if diff := (testutils.Problems{}).Diff(fileOptionConsistency.Lint(testFile)); diff != "" {
			t.Errorf(diff)
		}
	})
}

func getOptions(fileopts map[string]string) *dpb.FileOptions {
	opts := &dpb.FileOptions{
		CsharpNamespace: proto.String("Google.Example.V1"),
		GoPackage:       proto.String("github.com/googleapis/genproto/googleapis/example/v1;example"),
		JavaPackage:     proto.String("com.google.example.v1"),
		PhpNamespace:    proto.String("Google\\Example\\V1"),
		ObjcClassPrefix: proto.String("GEX"),
		RubyPackage:     proto.String("Google::Example::V1"),
		SwiftPrefix:     proto.String("Google_Example_V1"),
	}
	for k, v := range fileopts {
		field := strcase.UpperCamelCase(k)
		reflect.ValueOf(opts).Elem().FieldByName(field).Set(reflect.ValueOf(proto.String(v)))
	}
	return opts
}

type consistencyTest struct {
	name    string
	options map[string]string
}

func (ct *consistencyTest) getProblems(f *desc.FileDescriptor) testutils.Problems {
	p := testutils.Problems{}
	for k := range ct.options {
		p = append(p, lint.Problem{Message: k, Descriptor: f})
	}
	sort.Slice(p, func(i, j int) bool {
		return p[i].Message < p[j].Message
	})
	return p
}
