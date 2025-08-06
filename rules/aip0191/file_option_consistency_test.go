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
	"sort"
	"testing"

	"github.com/googleapis/api-linter/v2/rules/internal/testutils"
)

func TestFileOptionConsistency(t *testing.T) {
	controlSrc := `
		package google.example.v1;
		option csharp_namespace = "Google.Example.V1";
		option go_package = "github.com/googleapis/genproto/googleapis/example/v1;example";
		option java_package = "com.google.example.v1";
		option java_multiple_files = true;
		option php_namespace = "Google\\Example\\V1";
		option objc_class_prefix = "GEX";
		option ruby_package = "Google::Example::V1";
		option swift_prefix = "Google_Example_V1";
	`
	baseOptions := map[string]string{
		"csharp_namespace":       "Google.Example.V1",
		"go_package":             "github.com/googleapis/genproto/googleapis/example/v1;example",
		"java_package":           "com.google.example.v1",
		"java_multiple_files":    "true",
		"php_namespace":          "Google\\\\Example\\\\V1",
		"objc_class_prefix":      "GEX",
		"ruby_package":           "Google::Example::V1",
		"swift_prefix":           "Google_Example_V1",
		"php_metadata_namespace": "",
		"php_class_prefix":       "",
	}

	// Helper function to create a new map with modified options.
	// This avoids modifying the base map.
	withOptions := func(mods map[string]string) map[string]string {
		res := map[string]string{}
		for k, v := range baseOptions {
			res[k] = v
		}
		for k, v := range mods {
			if v == "" {
				delete(res, k)
			} else {
				res[k] = v
			}
		}
		return res
	}

	// Run our tests. Each test will make a "test file" that imports the
	// control file.
	for _, test := range []struct {
		name     string
		options  map[string]string
		problems testutils.Problems
	}{
		{"ValidSame", baseOptions, testutils.Problems{}},
		{"InvalidCsharpNamespaceMissing", withOptions(map[string]string{"csharp_namespace": ""}), testutils.Problems{{Message: `Option "csharp_namespace" should be consistent throughout the package.`}}},
		{"InvalidCsharpNamespaceMismatch", withOptions(map[string]string{"csharp_namespace": "Example.V1"}), testutils.Problems{{Message: `Option "csharp_namespace" should be consistent throughout the package.`}}},
		{"InvalidJavaPackageMissing", withOptions(map[string]string{"java_package": ""}), testutils.Problems{{Message: `Option "java_package" should be consistent throughout the package.`}}},
		{"InvalidJavaPackageMismatch", withOptions(map[string]string{"java_package": "com.example.v1"}), testutils.Problems{{Message: `Option "java_package" should be consistent throughout the package.`}}},
		{"InvalidPhpNamespaceMissing", withOptions(map[string]string{"php_namespace": ""}), testutils.Problems{{Message: `Option "php_namespace" should be consistent throughout the package.`}}},
		{"InvalidPhpNamespaceMismatch", withOptions(map[string]string{"php_namespace": "Example\\\\V1"}), testutils.Problems{{Message: `Option "php_namespace" should be consistent throughout the package.`}}},
		{"InvalidPhpMetadataNamespace", withOptions(map[string]string{"php_metadata_namespace": "Example\\\\V1"}), testutils.Problems{{Message: `Option "php_metadata_namespace" should be consistent throughout the package.`}}},
		{"InvalidPhpClassPrefix", withOptions(map[string]string{"php_class_prefix": "ExampleProto"}), testutils.Problems{{Message: `Option "php_class_prefix" should be consistent throughout the package.`}}},
		{"InvalidObjcClassPrefixMissing", withOptions(map[string]string{"objc_class_prefix": ""}), testutils.Problems{{Message: `Option "objc_class_prefix" should be consistent throughout the package.`}}},
		{"InvalidObjcClassPrefixMismatch", withOptions(map[string]string{"objc_class_prefix": "GEXV1"}), testutils.Problems{{Message: `Option "objc_class_prefix" should be consistent throughout the package.`}}},
		{"InvalidRubyPackageMissing", withOptions(map[string]string{"ruby_package": ""}), testutils.Problems{{Message: `Option "ruby_package" should be consistent throughout the package.`}}},
		{"InvalidRubyPackageMismatch", withOptions(map[string]string{"ruby_package": "Example::V1"}), testutils.Problems{{Message: `Option "ruby_package" should be consistent throughout the package.`}}},
		{"InvalidSwiftPrefixMissing", withOptions(map[string]string{"swift_prefix": ""}), testutils.Problems{{Message: `Option "swift_prefix" should be consistent throughout the package.`}}},
		{"InvalidSwiftPrefixMismatch", withOptions(map[string]string{"swift_prefix": "ExampleProto"}), testutils.Problems{{Message: `Option "swift_prefix" should be consistent throughout the package.`}}},
		{"InvalidLotsOfReasons", withOptions(map[string]string{
			"csharp_namespace": "Example.V1",
			"go_package":       "google.golang.org/googleapis/genproto/googleapis/example/v1;example",
			"java_package":     "com.example.v1",
			"ruby_package":     "Example::V1",
			"swift_prefix":     "",
		}), testutils.Problems{
			{Message: `Option "csharp_namespace" should be consistent throughout the package.`},
			{Message: `Option "go_package" should be consistent throughout the package.`},
			{Message: `Option "java_package" should be consistent throughout the package.`},
			{Message: `Option "ruby_package" should be consistent throughout the package.`},
			{Message: `Option "swift_prefix" should be consistent throughout the package.`},
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			// Build our test file, which imports the control file.
			testSrc := `
				import "control.proto";
				package google.example.v1;
				{{range $key, $value := .Options}}
				option {{$key}} = "{{$value}}";
				{{end}}
			`
			// The template is basic, and does not handle bools.
			// We address this by replacing the quoted "true" with the raw bool.
			// This is a bit of a hack, but it works for this test.
			if val, ok := test.options["java_multiple_files"]; ok && val == "true" {
				delete(test.options, "java_multiple_files")
				testSrc += "\noption java_multiple_files = true;"
			}

			files := testutils.ParseProto3Tmpls(t, map[string]string{
				"control.proto": controlSrc,
				"test.proto":    testSrc,
			}, map[string]interface{}{
				"Options": test.options,
			})
			testFile := files["test.proto"]
			problems := fileOptionConsistency.Lint(testFile)
			// We need to sort the problems to have a consistent diff.
			sort.Slice(problems, func(i, j int) bool {
				return problems[i].Message < problems[j].Message
			})
			if diff := test.problems.SetDescriptor(testFile).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}

	// Also ensure separate packages are ignored.
	t.Run("ValidDifferentPackages", func(t *testing.T) {
		controlSrc := `
			package google.example.v1;
			option csharp_namespace = "Google.Example.V1";
			option go_package = "github.com/googleapis/genproto/googleapis/example/v1;example";
			option java_package = "com.google.example.v1";
			option php_namespace = "Google\\Example\\V1";
			option objc_class_prefix = "GEX";
			option ruby_package = "Google::Example::V1";
			option swift_prefix = "Google_Example_V1";
		`
		testSrc := `
			import "control.proto";
			package google.different.v1;
			option java_package = "com.wrong";
		`
		files := testutils.ParseProto3Tmpls(t, map[string]string{
			"control.proto": controlSrc,
			"test.proto":    testSrc,
		}, nil)
		testFile := files["test.proto"]
		if diff := (testutils.Problems{}).Diff(fileOptionConsistency.Lint(testFile)); diff != "" {
			t.Error(diff)
		}
	})
}
