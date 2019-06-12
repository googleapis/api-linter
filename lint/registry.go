package lint

import (
	"fmt"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
)

// MakeRegistryFromAllFiles creates a *protoregistry.Files with all dependencies resolved, provided that
// any import in any FileDescriptorProto in files is contained in files.
//
// In other words, if for any i, files[i] imports "a.proto", then the FileDescriptorProto for "a.proto"
// must also be present in files.
func MakeRegistryFromAllFiles(descs []*descriptorpb.FileDescriptorProto) (*protoregistry.Files, error) {
	fileMap, err := makeFileMap(descs)
	if err != nil {
		return nil, err
	}

	f := files{
		reg:      new(protoregistry.Files),
		filesSet: fileMap,
	}

	for _, desc := range descs {
		if err := f.register(desc.GetName()); err != nil {
			return nil, err
		}
	}

	return f.reg, nil
}

type files struct {
	reg *protoregistry.Files

	filesSet map[string]*entry
}

type entry struct {
	desc        protoreflect.FileDescriptor
	registered  bool
	registering bool
}

func (f *files) register(path string) error {
	e, ok := f.filesSet[path]

	if !ok {
		return fmt.Errorf("%q not found in provided FileDescriptorProtos", path)
	}

	if e.registered {
		return nil
	}

	if e.registering {
		return fmt.Errorf("cyclic dependency found on import of %q", path)
	}

	e.registering = true

	for i := 0; i < e.desc.Imports().Len(); i++ {
		dep := e.desc.Imports().Get(i).Path()

		err := f.register(dep)

		if err != nil {
			return err
		}
	}

	if err := f.reg.Register(e.desc); err != nil {
		return err
	}

	e.registering = false
	e.registered = true

	return nil
}

func makeFileMap(files []*descriptorpb.FileDescriptorProto) (map[string]*entry, error) {
	fileMap := make(map[string]*entry, len(files))

	for _, f := range files {
		fd, err := protodesc.NewFile(f, nil)

		if err != nil {
			return nil, err
		}

		fileMap[f.GetName()] = &entry{desc: fd}
	}

	return fileMap, nil
}
