// Package testutil provides helpers for testing the linter and rules.
package testutil

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/googleapis/api-linter/parser/protoc"
	"google.golang.org/protobuf/types/descriptorpb"
)

// FileDescriptorSpec defines a specification for generating a FileDescriptorProto
type FileDescriptorSpec struct {
	// Filename is the output of the returned FileDescriptorProto.GetName().
	Filename string
	// Template defines a text/template to use for the proto source.
	Template string
	// Data is plugged into the template to create the full source code.
	Data interface{}
	// Deps are any additional FileDescriptorProtos that the protocol compiler will need for the source.
	Deps []*descriptorpb.FileDescriptorProto
	// AdditionalProtoPaths are any additional proto_paths that the protocol compiler will need for the source.
	AdditionalProtoPaths []string
}

// MustCreateFileDescriptorProto creates a *descriptorpb.FileDescriptorProto from a string template and data.
// TODO: refactor it to `FileDescriptorProto(t *testing.T, template string, data interface{}, options ...protoc.Option) *descriptorpb.FileDescriptorProto`
func MustCreateFileDescriptorProto(spec FileDescriptorSpec) *descriptorpb.FileDescriptorProto {
	tmpDir, err := ioutil.TempDir("", "testing")
	if err != nil {
		log.Fatalln(err)
	}
	defer os.RemoveAll(tmpDir)

	filename := "test.proto"
	if spec.Filename != "" {
		filename = spec.Filename
	}

	f, err := os.Create(filepath.Join(tmpDir, filename))
	if err != nil {
		log.Fatalln(err)
	}

	if err := template.Must(template.New("").Parse(spec.Template)).Execute(f, spec.Data); err != nil {
		log.Fatalf("Error executing template %v", err)
	}
	f.Close()

	p := protoc.New(
		protoc.AddProtoPath(spec.AdditionalProtoPaths...),
		protoc.AddProtoPath(tmpDir),
		protoc.AddDescriptorSetIn(spec.Deps...),
	)

	files, err := p.Parse(f.Name())
	if err != nil {
		log.Fatalln(err)
		return nil
	}

	return getFile(files, filename)
}

func getFile(files *descriptorpb.FileDescriptorSet, name string) *descriptorpb.FileDescriptorProto {
	for _, file := range files.GetFile() {
		if file.GetName() == name {
			return file
		}
	}
	return nil
}
