package testutil

import (
	"testing"
)

func TestDescriptorFromProtoSource_ApiCommonProtoImport(t *testing.T) {
	desc := MustCreateFileDescriptorProtoFromTemplate(FileDescriptorSpec{
		Template: `syntax = "proto3";

import "google/api/auth.proto";

message Foo {
	google.api.Authentication foo = 1;
}`,
	})

	if len(desc.GetDependency()) != 1 {
		t.Fatalf("desc.GetDependency()=%d; want 1", len(desc.GetDependency()))
	}

	if want := "google/api/auth.proto"; desc.GetDependency()[0] != want {
		t.Fatalf("desc.GetDependency()[0] = %q; want %q", desc.GetDependency()[0], want)
	}
}
