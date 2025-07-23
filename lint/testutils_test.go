package lint

import (
	"testing"

	"context"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/bufbuild/protocompile"
)

func buildFile(t *testing.T, content string) protoreflect.FileDescriptor {
	t.Helper()

	protoPath := "test.proto"

	// Create a map to hold our in-memory proto content, keyed by its path.
	inMemorySources := map[string]string{
		protoPath: content,
	}

	// Create a SourceResolver that will use our in-memory map to access file content.
	resolver := &protocompile.SourceResolver{
		Accessor: protocompile.SourceAccessorFromMap(inMemorySources),
		// No ImportPaths are needed here since we're only resolving "test.proto"
		// which is directly provided by our Accessor.
	}

	// Create a new protocompile.Compiler and set its Resolver.
	compiler := protocompile.Compiler{
		Resolver: resolver,
		// Explicitly set SourceInfoMode to ensure comments and source positions are included.
		SourceInfoMode: protocompile.SourceInfoStandard,
	}

	// Call Compile with the string path of the proto file.
	// The compiler will use the configured SourceResolver to retrieve the content for "test.proto".
	fds, err := compiler.Compile(context.Background(), protoPath)
	if err != nil {
		t.Fatalf("Failed to compile proto content: %v", err)
	}

	// The Compile method returns a slice of FileDescriptors.
	// Since we only passed one input file, we expect one output descriptor.
	if len(fds) == 0 {
		t.Fatalf("No file descriptors returned after compilation")
	}

	// The protocompile's protoreflect.FileDescriptor type is an alias for
	// google.golang.org/protobuf/reflect/protoreflect.FileDescriptor, so it's
	// directly compatible.
	return fds[0]
}
