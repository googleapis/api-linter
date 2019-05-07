// The protoc-gen-api_linter binary is a protoc plugin that checks API definition in .proto files.
package main

import (
	"fmt"
	"os"

	"github.com/googleapis/api-linter/cmd/protoc-gen-api_linter/linter"
	"github.com/googleapis/api-linter/cmd/protoc-gen-api_linter/protogen"
)

func main() {
	if len(os.Args) > 1 {
		fmt.Fprintln(os.Stderr, "protoc-gen_api_linter: This program should be run by protoc, not directly!")
		fmt.Fprintln(os.Stderr, "Usage: protoc --api_linter_out=cfg_file=my_cfg_file,out_file=my_lint_output_file:. my_proto_file")
		os.Exit(1)
	}

	protogen.Run(linter.New())
}
