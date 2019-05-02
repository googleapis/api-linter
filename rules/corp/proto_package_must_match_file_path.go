package corp

import (
	"fmt"
	"github.com/golang/protobuf/v2/reflect/protoreflect"
	"github.com/jgeewax/api-linter/lint"
	"github.com/jgeewax/api-linter/rules/descriptor"
	"os"
	"regexp"
	"strings"
)

func protoPackageMustMatchFilePath() lint.Rule {
	return &descriptor.CallbackRule{
		RuleInfo: lint.RuleInfo{
			Name:         "corp::proto_package_must_match_file_path",
			Description:  "All protos must belong to a package that matches their file path (relative to google3/)",
			URI:          "http://go/ce-api-conformance-checks",
			RequestTypes: []lint.RequestType{lint.ProtoRequest},
		},
		Callback: descriptor.Callbacks{
			FileCallback: func(fileDescriptor protoreflect.FileDescriptor, source lint.DescriptorSource) (problems []lint.Problem, e error) {
				packageFromPath := pathToPackage(fileDescriptor.Path())

				if packageFromPath == "" {
					return []lint.Problem{{
						Message:    "Failed to process file path. Please make sure it is in a subdirectory of google3 and ends with .proto",
						Descriptor: fileDescriptor,
					}}, nil
				}

				if string(fileDescriptor.Package()) != packageFromPath {
					return []lint.Problem{{
						Message: fmt.Sprintf(
							"The file path %q must match the package path %q (with path separators replaced by dots",
							fileDescriptor.Path(), fileDescriptor.Package().Name()),
						Descriptor: fileDescriptor,
					}}, nil
				}

				return nil, nil
			},
		},
	}
}

var protoPathRegexp = regexp.MustCompile(`^/?(?:[^/]+/)*?google3/(.*)/[^/]+\.proto$`)

func pathToPackage(path string) string {
	extractedPath := protoPathRegexp.FindStringSubmatch(path)

	if len(extractedPath) < 2 {
		return ""
	}

	return strings.Replace(extractedPath[1], string(os.PathSeparator), ".", -1)
}
