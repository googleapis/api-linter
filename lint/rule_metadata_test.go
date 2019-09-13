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
	"testing"
)

func TestGetName(t *testing.T) {
	want := NewRuleName("foo")
	rm := RuleMetadata{Name: want}
	if got := rm.GetName(); got != want {
		t.Errorf("GetName returned %q, expected %q.", got, want)
	}
}

func TestGetDescription(t *testing.T) {
	want := "foo bar baz"
	rm := RuleMetadata{Description: want}
	if got := rm.GetDescription(); got != want {
		t.Errorf("GetDescription returned %q, expected %q.", got, want)
	}
}

func TestGetURI(t *testing.T) {
	want := "https://aip.dev/1"
	rm := RuleMetadata{URI: want}
	if got := rm.GetURI(); got != want {
		t.Errorf("GetURI returned %q, expected %q.", got, want)
	}
}
