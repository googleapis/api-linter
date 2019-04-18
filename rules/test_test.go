package rules

import (
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
		err := closeAndRemoveFile(f)

		if err != nil {
			log.Fatalf("Error removing proto file: %v", err)
		}
	}()

	_, err = f.Write([]byte(source))

	if err != nil {
		return nil, err
	}

	descSetF, err := ioutil.TempFile(tmpDir, "descset*")

	if err != nil {
		return nil, err
	}

	defer func() {
		err := closeAndRemoveFile(descSetF)

		if err != nil {
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

	cmd.Stderr = os.Stderr

	err = cmd.Run()

	if err != nil {
		return nil, err
	}

	descSet, err := ioutil.ReadFile(descSetF.Name())

	if err != nil {
		return nil, err
	}

	protoset := &descriptorpb.FileDescriptorSet{}
	if err := proto.Unmarshal(descSet, protoset); err != nil {
		log.Fatalf("Unable to parse %T from source: %v.", protoset, err)
	}

	return protoset.GetFile()[0], nil
}

func closeAndRemoveFile(f *os.File) error {
	err := f.Close()

	if err != nil {
		return err
	}

	err = os.Remove(f.Name())

	if err != nil {
		return err
	}

	return nil
}
