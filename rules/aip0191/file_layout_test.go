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
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
	"github.com/jhump/protoreflect/desc/builder"
)

func TestFileLayout(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		f := testutils.ParseProto3String(t, `
			service Library {}
			message Book {}
			enum Format {
				FORMAT_UNSPECIFIED = 0;
			}
		`)
		if problems := fileLayout.Lint(f); len(problems) > 0 {
			t.Errorf("%v", problems)
		}
	})

	t.Run("InvalidServiceAfterMessage", func(t *testing.T) {
		f := testutils.ParseProto3String(t, `
			message Book {}
			service Library {}
		`)
		want := testutils.Problems{{Descriptor: f.Services()[0]}}
		if diff := want.Diff(fileLayout.Lint(f)); diff != "" {
			t.Error(diff)
		}
	})

	t.Run("InvalidEnumBeforeMessage", func(t *testing.T) {
		f := testutils.ParseProto3String(t, `
			enum Format {
				FORMAT_UNSPECIFIED = 0;
			}
			message Book {}
		`)
		want := testutils.Problems{{Descriptor: f.Enums()[0]}}
		if diff := want.Diff(fileLayout.Lint(f)); diff != "" {
			t.Error(diff)
		}
	})

	t.Run("NoSourceInfo", func(t *testing.T) {
		f, err := builder.NewFile("test.proto").AddMessage(
			builder.NewMessage("Book"),
		).AddService(
			builder.NewService("Library"),
		).Build()
		if err != nil {
			t.Fatalf("Failed to build proto.")
		}
		if problems := fileLayout.Lint(f); len(problems) > 0 {
			t.Errorf("%v", problems)
		}
	})
}
