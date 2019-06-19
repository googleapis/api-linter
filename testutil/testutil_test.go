package testutil

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"
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
