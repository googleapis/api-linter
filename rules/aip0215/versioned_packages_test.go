// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aip0215

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestVersionedPackages(t *testing.T) {
	for _, test := range []struct {
		name        string
		PackageStmt string
		problems    testutils.Problems
	}{
		{"Stable", "package foo.bar.v1;", testutils.Problems{}},
		{"StableBignum", "package foo.bar.v999;", testutils.Problems{}},
		{"Alpha", "package foo.bar.v1alpha;", testutils.Problems{}},
		{"Alpha1", "package foo.bar.v1alpha1;", testutils.Problems{}},
		{"Alpha12", "package foo.bar.v1alpha12;", testutils.Problems{}},
		{"Beta", "package foo.bar.v1beta;", testutils.Problems{}},
		{"Beta1", "package foo.bar.v1beta1;", testutils.Problems{}},
		{"Beta12", "package foo.bar.v1beta12;", testutils.Problems{}},
		{"P1Beta1", "package foo.bar.v1p1beta1;", testutils.Problems{}},
		{"P12Beta1", "package foo.bar.v1p12beta1;", testutils.Problems{}},
		{"P1Beta12", "package foo.bar.v1p1beta12;", testutils.Problems{}},
		{"P12Beta12", "package foo.bar.v1p12beta12;", testutils.Problems{}},
		{"EAP", "package foo.bar.v1eap;", testutils.Problems{}},
		{"EAP1", "package foo.bar.v1eap1;", testutils.Problems{}},
		{"Test", "package foo.bar.v1test;", testutils.Problems{}},
		{"Test1", "package foo.bar.v1test1;", testutils.Problems{}},
		{"Type", "package foo.bar.type;", testutils.Problems{}},
		{"Master", "package foo.bar.master;", testutils.Problems{}},
		{"VXMaster", "package foo.bar.v3master;", testutils.Problems{}},
		{"ValidSubpackage", "package foo.bar.v1.resources;", testutils.Problems{}},
		{"MasterSubpackage", "package foo.master.bar;", testutils.Problems{}},
		{"VXMasterSubpackage", "package foo.v3master.bar;", testutils.Problems{}},
		{"InvalidNoVersion", "package foo.bar;", testutils.Problems{{Message: "versioned packages"}}},
		{"IgnoredRPC", "package google.rpc.foobar;", testutils.Problems{}},
		{"IgnoredLRO", "package google.longrunning.foobar;", testutils.Problems{}},
		{"IgnoredAPI", "package google.api.foobar;", testutils.Problems{}},
		{"IgnoredNoPackage", "", testutils.Problems{}},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, "{{.PackageStmt}}", test)
			if diff := test.problems.SetDescriptor(f).Diff(versionedPackages.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
