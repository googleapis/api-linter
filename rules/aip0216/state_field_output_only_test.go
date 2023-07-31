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

package aip0216

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestStateFieldOutputOnly(t *testing.T) {
	tests := []struct {
		name          string
		FieldName     string
		FieldType     string
		FieldBehavior string
		problems      testutils.Problems
	}{
		// Accepted
		{"ValidState", "state", "State", "[(google.api.field_behavior) = OUTPUT_ONLY]", testutils.Problems{}},
		{"ValidOtherFieldName", "country", "State", "[(google.api.field_behavior) = OUTPUT_ONLY]", testutils.Problems{}},
		{"ValidStateSuffix", "state", "WritersBlockState", "[(google.api.field_behavior) = OUTPUT_ONLY]", testutils.Problems{}},

		// No Annotation
		{"InvalidState", "state", "State", "", testutils.Problems{{Message: "OUTPUT_ONLY"}}},
		{"InvalidWithSuffix", "state", "WritersBlockState", "", testutils.Problems{{Message: "OUTPUT_ONLY"}}},
		{"InvalidWithFieldName", "city", "State", "", testutils.Problems{{Message: "OUTPUT_ONLY"}}},

		// Ignored
		{"NotAStateField", "state", "StateOfDespair", "", testutils.Problems{}},
		{"NotAnEnum", "state", "StateOfState", "", testutils.Problems{}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
			import "google/api/field_behavior.proto";

			message Book {
				enum State {
					STATE_UNSPECIFIED = 0;
					ACTIVE = 1;
				}

				message StateOfState {
					string name = 1;
				}

				// state enums end with 'State'
				enum WritersBlockState {
					WRITERS_BLOCK_STATE_UNSPECIFIED = 0;
					BLOCKED = 1;
				}

				// not a state enum
				enum StateOfDespair {
					STATE_OF_DESPAIR_UNSPECIFIED = 0;
					NOTREALLY = 1;
				}

				string other_state = 1;
				{{.FieldType}} {{.FieldName}} = 2 {{.FieldBehavior}};
			}
		`, test)

			field := f.GetMessageTypes()[0].GetFields()[1]
			if diff := test.problems.SetDescriptor(field).Diff(stateFieldOutputOnly.Lint(f)); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
