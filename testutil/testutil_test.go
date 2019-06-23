package testutil

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"google.golang.org/protobuf/types/descriptorpb"
)

func TestDescriptorFromProtoSource_CustomProtoPaths(t *testing.T) {
	_, thisFilePath, _, _ := runtime.Caller(0)
	desc := MustCreateFileDescriptorProto(FileDescriptorSpec{
		AdditionalProtoPaths: []string{fmt.Sprintf("%s/%s", filepath.Dir(thisFilePath), "testdata")},
		Template: `syntax = "proto3";

import "sample.proto";

message Foo {
	testdata.Sample foo = 1;
}`,
	})

	if len(desc.GetDependency()) != 1 {
		t.Fatalf("desc.GetDependency()=%d; want 1", len(desc.GetDependency()))
	}

	if want := "sample.proto"; desc.GetDependency()[0] != want {
		t.Fatalf("desc.GetDependency()[0] = %q; want %q", desc.GetDependency()[0], want)
	}
}

func TestDescriptorFromProtoSource_CustomDeps(t *testing.T) {
	foo := MustCreateFileDescriptorProto(FileDescriptorSpec{
		Filename: "foo.proto",
		Template: `syntax = "proto3";

message Foo {
	string foo = 1;
}`,
	})

	bar := MustCreateFileDescriptorProto(FileDescriptorSpec{
		Filename: "bar.proto",
		Template: `syntax = "proto3";

import "foo.proto";

message Bar {
	Foo foo = 1;
}`,
		Deps: []*descriptorpb.FileDescriptorProto{foo},
	})

	if got, want := len(bar.GetDependency()), 1; got != want {
		t.Fatalf("len(bar.GetDependency()) = %d; want %d", got, want)
	}

	if got, want := bar.GetDependency()[0], "foo.proto"; got != want {
		t.Fatalf("bar.GetDependency()[0] = %s; want %s", got, want)
	}
}

func TestDescriptorFromProtoSource_CommonProtos(t *testing.T) {
	desc := MustCreateFileDescriptorProto(FileDescriptorSpec{
		Template: `
		syntax = "proto3";

		import "google/type/date.proto";

		message Foo {
			google.type.Date date = 1;
		}`,
	})

	if len(desc.GetDependency()) != 1 {
		t.Fatalf("desc.GetDependency()=%d; want 1", len(desc.GetDependency()))
	}

	if want := "google/type/date.proto"; desc.GetDependency()[0] != want {
		t.Fatalf("desc.GetDependency()[0] = %q; want %q", desc.GetDependency()[0], want)
	}
}
