package testutil

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"google.golang.org/protobuf/types/descriptorpb"
)

func TestDescriptorFromProtoSource_CustomProtoPaths(t *testing.T) {
	dirName, err := ioutil.TempDir("", "api_linter")
	if err != nil {
		t.Fatalf("Failed to create dependency directory: %s", err)
	}
	defer os.RemoveAll(dirName)

	fh, err := os.Create(filepath.Join(dirName, "sample.proto"))
	if err != nil {
		t.Fatalf("Failed to create sample.proto: %s", err)
	}

	_, err = fh.WriteString(`
		syntax = "proto3";
		package testdata;
		message Sample {
			string foo = 1;
		}
`)
	if err != nil {
		t.Fatalf("Failed to write proto source to sample.proto: %s", err)
	}

	desc := MustCreateFileDescriptorProto(t, FileDescriptorSpec{
		AdditionalProtoPaths: []string{dirName},
		Template: `
		syntax = "proto3";

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
	foo := MustCreateFileDescriptorProto(t, FileDescriptorSpec{
		Filename: "foo.proto",
		Template: `
		syntax = "proto3";

		message Foo {
			string foo = 1;
		}`,
	})

	bar := MustCreateFileDescriptorProto(t, FileDescriptorSpec{
		Filename: "bar.proto",
		Template: `
		syntax = "proto3";

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
