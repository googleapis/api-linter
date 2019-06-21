package rules

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/googleapis/api-linter/testutil"
)

// This function maps testdata directories to their real directories. The convoluted mechanism for
// doing so is an unfortunate consequence of maintaining compatibility when syncing internally.
var testdata = func(lib string) string {
	_, f, _, _ := runtime.Caller(0)
	thisDir := filepath.Dir(f) + string(os.PathSeparator)
	switch lib {
	case "api-common-protos":
		return thisDir + "testdata/api-common-protos"
	}
	return ""
}

func TestTestData_ApiCommonProtos(t *testing.T) {
	fd := testutil.MustCreateFileDescriptorProto(testutil.FileDescriptorSpec{
		Template: `syntax = "proto3";

import "google/api/auth.proto";

message Foo {
	google.api.Authentication foo = 1;
}
`,
		AdditionalProtoPaths: []string{testdata("api-common-protos")},
	})

	if got, want := len(fd.GetDependency()), 1; got != want {
		t.Fatalf("len(fd.GetDependency()) = %d; want %d", got, want)
	}

	if got, want := fd.GetDependency()[0], "google/api/auth.proto"; got != want {
		t.Fatalf("fd.GetDependency()[0] = %q; want %q", got, want)
	}
}
