package testutil

import (
	"strings"
	"testing"
)

func TestDescriptorFromProtoSource_ApiCommonProtoImport(t *testing.T) {
	desc, err := descriptorProtoFromSource(strings.NewReader(`syntax = "proto3";

import "google/api/auth.proto";

message Foo {
	google.api.Authentication foo = 1;
}`))

	if err != nil {
		t.Fatalf("Error creating descriptor: %s.", err)
	}

	if len(desc.GetDependency()) != 1 {
		t.Fatalf("desc.GetDependency()=%d; want 1", len(desc.GetDependency()))
	}

	if want := "google/api/auth.proto"; desc.GetDependency()[0] != want {
		t.Fatalf("desc.GetDependency()[0] = %q; want %q", desc.GetDependency()[0], want)
	}
}
