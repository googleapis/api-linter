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

package testutils

import (
	"testing"

	"github.com/jhump/protoreflect/desc"
)

func TestParseProtoString(t *testing.T) {
	fd := ParseProtoString(`
		syntax = "proto3";

		message Foo {
			int32 bar = 1;
			int64 baz = 2;
		}

		message Spam {
			string eggs = 2;
		}
	`)
	if !fd.IsProto3() {
		t.Errorf("Expected a proto3 file descriptor.")
	}
	tests := []struct {
		name       string
		descriptor desc.Descriptor
	}{
		{"Foo", fd.GetMessageTypes()[0]},
		{"bar", fd.GetMessageTypes()[0].GetFields()[0]},
		{"baz", fd.GetMessageTypes()[0].GetFields()[1]},
		{"Spam", fd.GetMessageTypes()[1]},
		{"eggs", fd.GetMessageTypes()[1].GetFields()[0]},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := test.descriptor.GetName(), test.name; got != want {
				t.Errorf("Got %q, expected %q.", got, want)
			}
		})
	}
}

func TestParseProtoStringError(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected a panic!")
		}
	}()
	ParseProtoString(`
		syntax = "proto3";
		message Foo {}
		The quick brown fox jumped over the lazy dogs.
	`)
}
