package testutils

import (
	"testing"

	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestCompileDescriptorProto(t *testing.T) {
	Compile(t, `
		import "google/protobuf/descriptor.proto";

		message Foo {
			google.protobuf.FileOptions options = 1;
		}
	`, nil)
}

func TestCompile_MessageAndField(t *testing.T) {
	file := Compile(t, `
		message MyMessage {
			string my_field = 1;
		}
	`, nil)

	// Verify the file descriptor
	if file == nil {
		t.Fatalf("Compile returned a nil FileDescriptor.")
	}
	if file.Messages().Len() != 1 {
		t.Fatalf("Expected 1 message, got %d", file.Messages().Len())
	}

	// Verify the message descriptor
	msg := file.Messages().ByName("MyMessage")
	if msg == nil {
		t.Fatalf("Expected message 'MyMessage' not found.")
	}
	if string(msg.FullName()) != "MyMessage" {
		t.Fatalf("Expected full name 'MyMessage', got '%s'", msg.FullName())
	}
	if msg.Fields().Len() != 1 {
		t.Fatalf("Expected 1 field in MyMessage, got %d", msg.Fields().Len())
	}

	// Verify the field descriptor
	field := msg.Fields().ByName("my_field")
	if field == nil {
		t.Fatalf("Expected field 'my_field' not found.")
	}
	if string(field.Name()) != "my_field" {
		t.Fatalf("Expected field name 'my_field', got '%s'", field.Name())
	}
	if field.Kind() != protoreflect.StringKind {
		t.Fatalf("Expected field kind StringKind, got %s", field.Kind())
	}
}

func TestCompile_FieldBehaviorMessage(t *testing.T) {
	Compile(t, `
		import "google/api/field_behavior.proto";

		message Message {
			string foo = 1 [(google.api.field_behavior) = REQUIRED];
		}
	`, nil)
}
