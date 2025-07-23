package lint

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

func buildFile(t *testing.T, content string) protoreflect.FileDescriptor {
	t.Helper()
	// Create a temporary directory.
	dir, err := ioutil.TempDir("", "proto-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	// Write the proto content to a file.
	tmpFN := filepath.Join(dir, "test.proto")
	if err := ioutil.WriteFile(tmpFN, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	// Run protoc to generate the file descriptor set.
	fdsFile := filepath.Join(dir, "fds.pb")
	cmd := exec.Command("protoc", "-o", fdsFile, "--include_imports", "--include_source_info", "test.proto")
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("protoc failed: %s\n%s", err, output)
	}

	// Read the file descriptor set.
	fdsBin, err := ioutil.ReadFile(fdsFile)
	if err != nil {
		t.Fatalf("Failed to read fds file: %v", err)
	}
	fds := &descriptorpb.FileDescriptorSet{}
	if err := proto.Unmarshal(fdsBin, fds); err != nil {
		t.Fatalf("Failed to unmarshal fds: %v", err)
	}

	// Create the file descriptor.
	files, err := protodesc.NewFiles(fds)
	if err != nil {
		t.Fatalf("Failed to create file descriptor: %v", err)
	}

	fd, err := files.FindFileByPath("test.proto")
	if err != nil {
		t.Fatalf("Failed to find file descriptor: %v", err)
	}
	return fd
}
