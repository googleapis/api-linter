// Package testutil provides helpers for testing the linter and rules.
package testutil

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"text/template"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

var protocPath = func() string {
	return "protoc"
}

// MustCreateFileDescriptorProtoFromTemplate creates a *descriptorpb.FileDescriptorProto from a string template and data.
func MustCreateFileDescriptorProtoFromTemplate(filename, srcTmpl string, data interface{}, deps []*descriptorpb.FileDescriptorProto) *descriptorpb.FileDescriptorProto {
	tmpl := template.Must(template.New("test").Parse(srcTmpl))
	b := new(bytes.Buffer)
	if err := tmpl.Execute(b, data); err != nil {
		log.Fatalf("Error executing template %v", err)
	}

	return mustCreateDescriptorProtoFromSource(filename, b, deps)
}

func mustCreateDescriptorProtoFromSource(filename string, source io.Reader, deps []*descriptorpb.FileDescriptorProto) *descriptorpb.FileDescriptorProto {
	tmpDir := os.TempDir()

	f, err := ioutil.TempFile(tmpDir, "proto*")
	if err != nil {
		log.Fatalf("Failed creating temp proto source file: %s", err)
	}
	defer mustCloseAndRemoveFile(f)

	if _, err = io.Copy(f, source); err != nil {
		log.Fatalf("Failed to copy source to templ file: %s", err)
	}

	descSetF, err := ioutil.TempFile(tmpDir, "descset*")
	if err != nil {
		log.Fatalf("Failed to create temp descriptor set file: %s", err)
	}
	defer mustCloseAndRemoveFile(descSetF)

	args := []string{
		"--include_source_info",
		fmt.Sprintf("--proto_path=%s", tmpDir),
		fmt.Sprintf("--descriptor_set_out=%s", descSetF.Name()),
	}

	if len(deps) > 0 {
		descSetIn := mustCreateDescSetFileFromFileDescriptorProtos(deps)
		defer mustCloseAndRemoveFile(descSetIn)

		args = append(args, fmt.Sprintf("--descriptor_set_in=%s", descSetIn.Name()))
	}

	args = append(args, f.Name())

	cmd := exec.Command(protocPath(), args...)

	stderr := new(bytes.Buffer)
	cmd.Stderr = stderr

	if err = cmd.Run(); err != nil {
		log.Fatalf("protoc failed with %v and Stderr %q", err, stderr.String())
	}

	descSet, err := ioutil.ReadFile(descSetF.Name())
	if err != nil {
		log.Fatalf("Failed to read descriptor set file: %s", err)
	}

	protoset := &descriptorpb.FileDescriptorSet{}
	if err := proto.Unmarshal(descSet, protoset); err != nil {
		log.Fatalf("Failed to unmarshal descriptor set file: %s", err)
	}

	if len(protoset.GetFile()) == 0 {
		log.Fatalf("protoset.GetFile() returns empty list")
	}

	protoset.GetFile()[0].Name = &filename

	return protoset.GetFile()[0]
}

func mustCreateDescSetFileFromFileDescriptorProtos(descs []*descriptorpb.FileDescriptorProto) *os.File {
	if len(descs) == 0 {
		return nil
	}

	descSet := new(descriptorpb.FileDescriptorSet)
	descSet.File = descs

	rawDescSet, err := proto.Marshal(descSet)

	if err != nil {
		log.Fatalf("Failed to marshal descriptor set: %s", err)
	}

	descSetF, err := ioutil.TempFile(os.TempDir(), "descset*")

	if err != nil {
		log.Fatalf("Failed to make descriptor set file: %s", err)
	}

	if _, err := io.Copy(descSetF, bytes.NewReader(rawDescSet)); err != nil {
		mustCloseAndRemoveFile(descSetF)
		log.Fatalf("Failed to ")
	}

	return descSetF
}

func mustCloseAndRemoveFile(f *os.File) {
	if err := f.Close(); err != nil {
		log.Fatalf("Error closing proto file: %v", err)
	}

	if err := os.Remove(f.Name()); err != nil {
		log.Fatalf("Error removing proto file: %v", err)
	}
}
