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

package utils

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/googleapis/api-linter/rules/internal/testutils"
	apb "google.golang.org/genproto/googleapis/api/annotations"
)

func TestGetFieldBehavior(t *testing.T) {
	fd := testutils.ParseProto3String(t, `
		import "google/api/field_behavior.proto";

		message Book {
			string name = 1 [
				(google.api.field_behavior) = IMMUTABLE,
				(google.api.field_behavior) = OUTPUT_ONLY];

			string title = 2 [(google.api.field_behavior) = REQUIRED];

			string summary = 3;
		}
	`)
	msg := fd.GetMessageTypes()[0]
	tests := []struct {
		fieldName      string
		fieldBehaviors []apb.FieldBehavior
	}{
		{"name", []apb.FieldBehavior{apb.FieldBehavior_IMMUTABLE, apb.FieldBehavior_OUTPUT_ONLY}},
		{"title", []apb.FieldBehavior{apb.FieldBehavior_REQUIRED}},
		{"summary", nil},
	}
	for _, test := range tests {
		t.Run(test.fieldName, func(t *testing.T) {
			f := msg.FindFieldByName(test.fieldName)
			if diff := cmp.Diff(GetFieldBehavior(f), test.fieldBehaviors); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
