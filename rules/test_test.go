package rules

import (
	"bytes"
	"fmt"
	"github.com/golang/protobuf/v2/proto"
	descriptorpb "github.com/golang/protobuf/v2/types/descriptor"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func protoDescriptorProtoFromSource(source string) (*descriptorpb.FileDescriptorProto, error) {
	tmpDir := os.TempDir()

	f, err := ioutil.TempFile(tmpDir, "proto*")

	if err != nil {
		return nil, err
	}

	defer func() {
		if err := closeAndRemoveFile(f); err != nil {
			log.Fatalf("Error removing proto file: %v", err)
		}
	}()

	if _, err = f.WriteString(source); err != nil {
		return nil, err
	}

	descSetF, err := ioutil.TempFile(tmpDir, "descset*")

	if err != nil {
		return nil, err
	}

	defer func() {
		if err := closeAndRemoveFile(descSetF); err != nil {
			log.Fatalf("Error removing descriptor set file: %v", err)
		}
	}()

	cmd := exec.Command(
		"protoc",
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

func closeAndRemoveFile(f *os.File) error {
	if err := f.Close(); err != nil {
		return err
	}

	if err := os.Remove(f.Name()); err != nil {
		return err
	}

	return nil
}
