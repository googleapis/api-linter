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

func TestMethodRequestType(t *testing.T) {
	f := parse(t, `
		service Library {
		  rpc GetBook(GetBookRequest) returns (Book);
		}
		message GetBookRequest {}
		message Book {}
	`)
	loc := MethodRequestType(f.GetServices()[0].GetMethods()[0])
	// Three character span: line, start column, end column.
	if diff := cmp.Diff(loc.GetSpan(), []int32{3, 14, 28}); diff != "" {
		t.Errorf(diff)
	}
}

func TestMethodResponseType(t *testing.T) {
	f := parse(t, `
		service Library {
		  rpc GetBook(GetBookRequest) returns (Book);
		}
		message GetBookRequest {}
		message Book {}
	`)
	loc := MethodResponseType(f.GetServices()[0].GetMethods()[0])
	// Three character span: line, start column, end column.
	if diff := cmp.Diff(loc.GetSpan(), []int32{3, 39, 43}); diff != "" {
		t.Errorf(diff)
	}
}

func TestMethodHTTPRule(t *testing.T) {
	f := parse(t, `
		import "google/api/annotations.proto";
		service Library {
		  rpc GetBook(GetBookRequest) returns (Book) {
		    option (google.api.http) = {
		      get: "/v1/{name=publishers/*/books/*}"
		    };
		  }
		}
		message GetBookRequest{}
		message Book {}
	`)
	loc := MethodHTTPRule(f.GetServices()[0].GetMethods()[0])
	// Four character span: start line, start column, end line, end column.
	if diff := cmp.Diff(loc.GetSpan(), []int32{5, 4, 7, 6}); diff != "" {
		t.Errorf(diff)
	}
}

func TestMethodSignature(t *testing.T) {
	f := parse(t, `
		import "google/api/client.proto";
		service Library {
		  rpc GetBook(GetBookRequest) returns (Book) {
		    option (google.api.method_signature) = "name";
		    option (google.api.method_signature) = "name,read_mask";
		  }
		}
		message GetBookRequest{}
		message Book {}
	`)
	for _, test := range []struct {
		name  string
		index int
		want  []int32
	}{
		{"First", 0, []int32{5, 4, 50}},
		{"Second", 1, []int32{6, 4, 60}},
	} {
		loc := MethodSignature(f.GetServices()[0].GetMethods()[0], test.index)
		// Four character span: start line, start column, end line, end column.
		if diff := cmp.Diff(loc.GetSpan(), test.want); diff != "" {
			t.Errorf(diff)
		}
	}
}
