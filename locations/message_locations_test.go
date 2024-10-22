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
)

func TestMessageResource(t *testing.T) {
	f := parse(t, `
		import "google/api/resource.proto";
		message Book {
		  option (google.api.resource) = {
		    type: "library.googleapis.com/Book"
		    pattern: "publishers/{publisher}/books/{book}"
		  };
		}
	`)
	loc := MessageResource(f.GetMessageTypes()[0])
	if diff := cmp.Diff(loc.GetSpan(), []int32{4, 2, 7, 4}); diff != "" {
		t.Error(diff)
	}
}
