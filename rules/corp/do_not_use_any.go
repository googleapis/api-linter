package corp

import (
	"github.com/golang/protobuf/v2/reflect/protoreflect"
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/descriptor"
	"log"
)

var anyPath = "google/protobuf/any.proto"

func doNotUseAny() lint.Rule {
	return &descriptor.CallbackRule{
		RuleInfo: lint.RuleInfo{
			Name:         "corp::do_not_use_any",
			Description:  "The Any proto should not be used.",
			URI:          "http://go/ce-api-conformance-checks",
			RequestTypes: []lint.RequestType{lint.ProtoRequest},
		},
		Callback: descriptor.Callbacks{
			FileCallback: func(fileDescriptor protoreflect.FileDescriptor, source lint.DescriptorSource) (problems []lint.Problem, e error) {
				for i := 0; i < fileDescriptor.Imports().Len(); i++ {
					importDescriptor := fileDescriptor.Imports().Get(i)

					if string(importDescriptor.Path()) == anyPath {
						loc, err := source.DescriptorLocation(importDescriptor)

						if err != nil {
							log.Printf("Failed to get Location for Descriptor because %s.", err)
							return
						}

						problems = append(problems, lint.Problem{
							Message:    "Found use of Any proto.",
							Descriptor: importDescriptor,
							Location:   loc,
						})
					}
				}

				return
			},
		},
	}
}
