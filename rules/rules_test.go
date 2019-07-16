package rules

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/googleapis/api-linter/testutil"
)

// This function maps testdata directories to their real directories. The convoluted mechanism for
// doing so is an unfortunate consequence of maintaining compatibility when syncing internally.
var testdatadir = func(lib string) string {
	_, f, _, _ := runtime.Caller(0)
	testdataDir := filepath.Join(filepath.Dir(f), "testdata")
	switch lib {
	case "api-common-protos":
		return filepath.Join(testdataDir, "api-common-protos")
	}
	return ""
}

func TestTestData_ApiCommonProtos(t *testing.T) {
	fd := testutil.MustCreateFileDescriptorProto(t, testutil.FileDescriptorSpec{
		Filename: "test.proto",
		Template: `syntax = "proto3";

import "google/api/auth.proto";

message Foo {
	google.api.Authentication foo = 1;
}
`,
		AdditionalProtoPaths: []string{testdatadir("api-common-protos")},
	})

	if got, want := len(fd.GetDependency()), 1; got != want {
		t.Fatalf("len(fd.GetDependency()) = %d; want %d", got, want)
	}

	if got, want := fd.GetDependency()[0], "google/api/auth.proto"; got != want {
		t.Fatalf("fd.GetDependency()[0] = %q; want %q", got, want)
	}
}
