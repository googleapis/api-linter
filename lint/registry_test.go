package lint

import (
	"github.com/googleapis/api-linter/testutil"
	"google.golang.org/protobuf/types/descriptorpb"
	"testing"
)

func TestMakeRegistryFromAllFiles(t *testing.T) {
	barProto := testutil.MustCreateFileDescriptorProtoFromSpec(testutil.FileDescriptorSpec{
		Filename: "bar.proto",
		Template: `syntax = "proto3";

message Bar {
  string baz = 1;
}`,
	})

	fooProto := testutil.MustCreateFileDescriptorProtoFromSpec(testutil.FileDescriptorSpec{
		Filename: "foo.proto",
		Template: `syntax = "proto3";

import "bar.proto";

message Foo {
  Bar bar = 1;
}`,
		Deps: []*descriptorpb.FileDescriptorProto{barProto},
	})

	reg, err := makeRegistryFromAllFiles([]*descriptorpb.FileDescriptorProto{fooProto, barProto})

	if err != nil {
		t.Fatalf("makeRegistryFromAllFiles() returned error %q; want nil", err)
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
	if foo.Fields().Get(0).Message().IsPlaceholder() {
		t.Fatalf("foo.Fields().Get(0).IsPlaceholder()=true; want false")
	}
}

func TestMakeRegistryFromAllFiles_DirectAndIndirectDependencies(t *testing.T) {
	barProto := testutil.MustCreateFileDescriptorProtoFromSpec(testutil.FileDescriptorSpec{
		Filename: "bar.proto",
		Template: `syntax = "proto3";

message Bar {
  string baz = 1;
}`,
	})

	fooProto := testutil.MustCreateFileDescriptorProtoFromSpec(testutil.FileDescriptorSpec{
		Filename: "foo.proto",
		Template: `syntax = "proto3";

import "bar.proto";

message Foo {
  Bar bar = 1;
}`,
		Deps: []*descriptorpb.FileDescriptorProto{barProto},
	})

	bazProto := testutil.MustCreateFileDescriptorProtoFromSpec(testutil.FileDescriptorSpec{
		Filename: "baz.proto",
		Template: `syntax = "proto3";

import "bar.proto";
import "foo.proto";

message Baz {
	Foo foo = 1;
	Bar bar = 2;
}
`,
		Deps: []*descriptorpb.FileDescriptorProto{barProto, fooProto},
	})

	reg, err := makeRegistryFromAllFiles([]*descriptorpb.FileDescriptorProto{fooProto, barProto, bazProto})
	if err != nil {
		t.Fatalf("makeRegistryFromAllFiles() returned error %q; want nil", err)
	}

	foo, err := reg.FindMessageByName("Foo")
	if err != nil {
		t.Fatalf("reg.FindMessageByName(%q) returned error %q; want nil", "Foo", err)
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
	if foo.Fields().Get(0).Message().IsPlaceholder() {
		t.Fatalf("foo.Fields().Get(0).IsPlaceholder()=true; want false")
	}

	baz, err := reg.FindMessageByName("Baz")
	if err != nil {
		t.Fatalf("reg.FindMessageByName(%q) returned error %q; want nil", "Baz", err)
	}

	if baz.Fields().Len() != 2 {
		t.Fatalf("baz.Fields.Len()=%d; want 2", baz.Fields().Len())
	}

	if baz.Fields().Get(0).Message() == nil {
		t.Fatalf("baz.Fields().Get(0).Message() was nil")
	}
	if baz.Fields().Get(0).Message().Name() != "Foo" {
		t.Fatalf("baz.Fields().Get(0).Message().Name() = %q; want %q", baz.Fields().Get(1).Message().Name(), "Foo")
	}
	if baz.Fields().Get(0).Message().IsPlaceholder() {
		t.Fatalf("baz.Fields().Get(0).IsPlaceholder()=true; want false")
	}

	if baz.Fields().Get(1).Message() == nil {
		t.Fatalf("baz.Fields().Get(1).Message() was nil")
	}
	if baz.Fields().Get(1).Message().Name() != "Bar" {
		t.Fatalf("baz.Fields().Get(1).Message().Name() = %q; want %q", baz.Fields().Get(1).Message().Name(), "Bar")
	}
	if baz.Fields().Get(1).Message().IsPlaceholder() {
		t.Fatalf("baz.Fields().Get(1).IsPlaceholder()=true; want false")
	}
}

func TestMakeRegistryFromAllFiles_MissingImports(t *testing.T) {
	barProto := testutil.MustCreateFileDescriptorProtoFromSpec(testutil.FileDescriptorSpec{
		Filename: "bar.proto",
		Template: `syntax = "proto3";

message Bar {
  string baz = 1;
}`,
	})

	fooProto := testutil.MustCreateFileDescriptorProtoFromSpec(testutil.FileDescriptorSpec{
		Filename: "foo.proto",
		Template: `syntax = "proto3";

import "bar.proto";

message Foo {
  Bar bar = 1;
}`,
		Deps: []*descriptorpb.FileDescriptorProto{barProto},
	})

	reg, err := makeRegistryFromAllFiles([]*descriptorpb.FileDescriptorProto{fooProto})

	if err != nil {
		t.Fatalf("makeRegistryFromAllFiles() returned error %q; want nil", err)
	}

	foo, err := reg.FindMessageByName("Foo")
	if err != nil {
		t.Fatalf("Failed to find Foo message: %s.", err)
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
	if !foo.Fields().Get(0).Message().IsPlaceholder() {
		t.Fatalf("foo.Fields().Get(0).IsPlaceholder()=true; want false")
	}
}
