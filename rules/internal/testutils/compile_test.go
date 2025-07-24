package testutils

import (
	"testing"
)

func TestCompileDescriptorProto(t *testing.T) {
	Compile(t, `
		import "google/protobuf/descriptor.proto";

		message Foo {
			google.protobuf.FileOptions options = 1;
		}
	`, nil)
}

func TestCompile_FieldBehaviorMessage(t *testing.T) {
	Compile(t, `
		import "google/api/field_behavior.proto";

		message Message {
			string foo = 1 [(google.api.field_behavior) = REQUIRED];
		}
	`, nil)
}
