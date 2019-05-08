// The protoc-gen-api_linter binary is a protoc plugin
// that checks API definition in .proto files.
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/googleapis/api-linter/cmd/protoc-gen-api_linter/linter"
	"github.com/googleapis/api-linter/cmd/protoc-gen-api_linter/protogen"
	"github.com/googleapis/api-linter/lint"
	corerules "github.com/googleapis/api-linter/rules"
)

// Register default configuration.
func getConfigs() lint.RuntimeConfigs {
	return lint.RuntimeConfigs{
		lint.RuntimeConfig{
			IncludedPaths: []string{"**/*.proto"},
			RuleConfigs: map[string]lint.RuleConfig{
				"": {
					Status:   lint.Disabled,
					Category: lint.Warning,
				},
			},
		},
	}
}

// Register rules.
func getRules() []lint.Rule {
	var rules []lint.Rule
	rules = append(rules, corerules.Rules().All()...)
	return rules
}

func main() {
	if len(os.Args) > 1 {
		fmt.Fprintln(os.Stderr, "protoc-gen_api_linter: This program should be run by protoc, not directly!")
		fmt.Fprintln(os.Stderr, "Usage: protoc --api_linter_out=cfg_file=my_cfg_file,out_file=my_lint_output_file:. my_proto_file")
		os.Exit(1)
	}
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", filepath.Base(os.Args[0]), err)
		os.Exit(1)
	}
}

func run() error {
	return protogen.Run(linter.New(getRules(), getConfigs()))
}
