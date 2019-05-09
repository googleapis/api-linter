package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/golang/protobuf/proto"
	descriptorpb "github.com/golang/protobuf/v2/types/descriptor"
)

var protoc = "protoc"

type protocParser struct {
	importPath string
}

func (p *protocParser) ParseProto(filenames ...string) ([]*descriptorpb.FileDescriptorProto, error) {
	return compileProto(p.importPath, filenames...)
}

func compileProto(importPath string, filenames ...string) ([]*descriptorpb.FileDescriptorProto, error) {
	workdir, err := ioutil.TempDir("", "test")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(workdir)

	outfile := filepath.Join(workdir, "desc.proto")
	cmd := exec.Command(protoc, "--descriptor_set_out="+outfile)
	cmd.Args = append(cmd.Args, "--include_source_info")
	cmd.Args = append(cmd.Args, "-I="+importPath)
	cmd.Args = append(cmd.Args, filenames...)
	if out, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("Running %s: %s, %v", strings.Join(cmd.Args, " "), string(out), err)
	}

	b, err := ioutil.ReadFile(outfile)
	if err != nil {
		return nil, err
	}
	v := &descriptorpb.FileDescriptorSet{}
	if err := proto.Unmarshal(b, v); err != nil {
		return nil, err
	}
	return v.GetFile(), nil
}
