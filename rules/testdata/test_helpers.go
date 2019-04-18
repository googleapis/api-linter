// Package testdata provides testing helpers and data for rules.
package testdata

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/golang/protobuf/v2/proto"
	descriptorpb "github.com/golang/protobuf/v2/types/descriptor"
	"github.com/jgeewax/api-linter/lint"
)

var protoc = "protoc"

func descriptorProtoFromSource(source []byte) (*descriptorpb.FileDescriptorProto, error) {
	tmpDir := os.TempDir()

	f, err := ioutil.TempFile(tmpDir, "proto*")

	if err != nil {
		return nil, err
	}
	defer mustCloseAndRemoveFile(f)

	if _, err = f.Write(source); err != nil {
		return nil, err
	}

	descSetF, err := ioutil.TempFile(tmpDir, "descset*")

	if err != nil {
		return nil, err
	}
	defer mustCloseAndRemoveFile(descSetF)

	cmd := exec.Command(
		protoc,
		"--include_source_info",
		fmt.Sprintf("--proto_path=%s", tmpDir),
		fmt.Sprintf("--descriptor_set_out=%s", descSetF.Name()),
		f.Name(),
	)

	var stdErrBuf bytes.Buffer

	cmd.Stderr = &stdErrBuf

	if err = cmd.Run(); err != nil {
		return nil, fmt.Errorf("protoc failed with %v and Stderr %q", err, stdErrBuf.String())
	}

	descSet, err := ioutil.ReadFile(descSetF.Name())

	if err != nil {
		return nil, err
	}

	protoset := &descriptorpb.FileDescriptorSet{}
	if err := proto.Unmarshal(descSet, protoset); err != nil {
		return nil, err
	}

	if len(protoset.GetFile()) == 0 {
		return nil, fmt.Errorf("protoset file list was empty")
	}

	return protoset.GetFile()[0], nil
}

func mustCloseAndRemoveFile(f *os.File) {
	if err := f.Close(); err != nil {
		log.Fatalf("Error closing proto file: %v", err)
	}

	if err := os.Remove(f.Name()); err != nil {
		log.Fatalf("Error removing proto file: %v", err)
	}
}

// MustCreateTemplate creates a template with name "test" from
// the provided template string.
func MustCreateTemplate(tmpl string) *template.Template {
	return template.Must(template.New("test").Parse(tmpl))
}

// MustCreateRequestFromTemplate creates a lint.Request from the provided template and test data.
func MustCreateRequestFromTemplate(tmpl *template.Template, testData interface{}) lint.Request {
	var b bytes.Buffer
	if err := tmpl.Execute(&b, testData); err != nil {
		log.Fatalf("Error executing template %v", err)
	}
	pd, err := descriptorProtoFromSource(b.Bytes())
	if err != nil {
		log.Fatalf("Error generating proto descriptor: %v", err)
	}
	req, err := lint.NewProtoFileRequest(pd)
	if err != nil {
		log.Fatalf("Error creating proto file request: %v", err)
	}
	return req
}
