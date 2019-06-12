package lint

import (
	"github.com/googleapis/api-linter/rules/testutil"
	"google.golang.org/protobuf/types/descriptorpb"
	"testing"
)

func TestMakeRegistryFromAllFiles(t *testing.T) {
	barProto := testutil.MustCreateFileDescriptorProtoFromTemplate("bar.proto", `syntax = "proto3";

message Bar {
  string baz = 1;
}`, nil, nil)

	fooProto := testutil.MustCreateFileDescriptorProtoFromTemplate("foo.proto", `syntax = "proto3";

import "bar.proto";

message Foo {
  Bar bar = 1;
}`, nil, []*descriptorpb.FileDescriptorProto{barProto})

	reg, err := MakeRegistryFromAllFiles([]*descriptorpb.FileDescriptorProto{fooProto, barProto})

	if err != nil {
		t.Fatalf("MakeRegistryFromAllFiles() returned error %q; want nil", err)
	}

	foo, err := reg.FindMessageByName("Foo")

	if err != nil {
		t.Fatalf("reg.FindMessageByName(%q) returned error %q; want nil", fooProto.GetName(), err)
	}

	if foo.Fields().Len() != 1 {
		t.Fatalf("foo.Fields().Len()=%d; want 1", foo.Fields().Len())
	}

	if foo.Fields().Get(0).Message() == nil {
		t.Fatalf("foo.Fields().Get(0).Message() was nil")
	}

	if foo.Fields().Get(0).Message().Name() != "Bar" {
		t.Fatalf("foo.Fields().Get(0).Message().Name() = %q; want %q", foo.Fields().Get(0).Message().Name(), "Bar")
	}
}

func TestMakeRegistryFromAllFiles_MissingImports(t *testing.T) {
	barProto := testutil.MustCreateFileDescriptorProtoFromTemplate("bar.proto", `syntax = "proto3";

message Bar {
  string baz = 1;
}`, nil, nil)

	fooProto := testutil.MustCreateFileDescriptorProtoFromTemplate("foo.proto", `syntax = "proto3";

import "bar.proto";

message Foo {
  Bar bar = 1;
}`, nil, []*descriptorpb.FileDescriptorProto{barProto})

	_, err := MakeRegistryFromAllFiles([]*descriptorpb.FileDescriptorProto{fooProto})

	if err == nil {
		t.Fatalf("MakeRegistryFromAllFiles() returned nil error, but there were missing imports.")
	}
}
