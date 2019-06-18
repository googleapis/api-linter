package lint

import (
	"fmt"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
)

// makeRegistryFromAllFiles creates a *protoregistry.Files with all dependencies present in descs resolved.
// Missing dependencies will be filled with placeholders. If there are cyclic dependencies, an error will
// be returned.
func makeRegistryFromAllFiles(descs []*descriptorpb.FileDescriptorProto) (*protoregistry.Files, error) {
	filesSet, err := makeFilesSet(descs)
	if err != nil {
		return nil, err
	}

	f := files{
		reg:      new(protoregistry.Files),
		filesSet: filesSet,
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
	descProto   *descriptorpb.FileDescriptorProto
	desc        protoreflect.FileDescriptor
	registering bool
}

func (f *files) register(path string) error {
	e, ok := f.filesSet[path]

	// if it's already been registered (from another import) do nothing. if it doesn't exist,
	// also do nothing, and a placeholder will be filled in for those that import it
	if !ok || e.desc != nil {
		return nil
	}

	if e.registering {
		return fmt.Errorf("cyclic dependency found on import of %q", path)
	}

	e.registering = true
	defer func() {
		e.registering = false
	}()

	for _, dep := range e.descProto.GetDependency() {
		if err := f.register(dep); err != nil {
			return err
		}
	}

	var err error
	if e.desc, err = protodesc.NewFile(e.descProto, f.reg); err != nil {
		return err
	}
	if err = f.reg.Register(e.desc); err != nil {
		return err
	}

	return nil
}

func makeFilesSet(files []*descriptorpb.FileDescriptorProto) (map[string]*entry, error) {
	filesSet := make(map[string]*entry, len(files))

	for _, f := range files {
		filesSet[f.GetName()] = &entry{descProto: f}
	}

	return filesSet, nil
}
