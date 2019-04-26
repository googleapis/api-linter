package corp

import (
	"fmt"
	"github.com/golang/protobuf/v2/reflect/protoreflect"
	"github.com/jgeewax/api-linter/lint"
	"github.com/jgeewax/api-linter/rules/descriptor"
	"os"
	"regexp"
)

var validVersionStrings = regexp.MustCompile(fmt.Sprintf(
	"%c(alpha|beta|v[\\d]+)%c",
	os.PathSeparator, os.PathSeparator,
))

func protoFilesMustIncludeVersion() lint.Rule {
	return &descriptor.CallbackRule{
		RuleInfo: lint.RuleInfo{
			Name:         "corp::proto_files_must_include_version",
			Description:  "All proto files must include a version number in their file path.",
			URI:          "http://go/ce-api-conformance-checks",
			RequestTypes: []lint.RequestType{lint.ProtoRequest},
		},
		Callback: descriptor.Callbacks{
			FileCallback: func(fileDescriptor protoreflect.FileDescriptor, source lint.DescriptorSource) (problems []lint.Problem, e error) {
				if !validVersionStrings.Match([]byte(fileDescriptor.Path())) {
					return []lint.Problem{{
						Message:    fmt.Sprintf("The file path %q must include a version specifier", fileDescriptor.Path()),
						Descriptor: fileDescriptor,
					}}, nil
				}

				return nil, nil
			},
		},
	}
}
