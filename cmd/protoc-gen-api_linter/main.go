// The protoc-gen-api_linter binary is a protoc plugin that checks API definition in .proto files.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	descriptorpb "github.com/golang/protobuf/v2/types/descriptor"
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules"
	"gopkg.in/yaml.v2"
)

func main() {
	var (
		flags     flag.FlagSet
		cfgFile   = flags.String("cfg_file", "", "Google API Linter configuration file")
		outFile   = flags.String("out_file", "", "Google API Linter output file")
		outFormat = flags.String("out_format", "", "Google API Linter output format")
		opts      = &Options{
			ParamFunc: flags.Set,
		}
	)

	if len(os.Args) > 1 {
		fmt.Fprintln(os.Stderr, "protoc-gen_api_linter: This program should be run by protoc, not directly!")
		fmt.Fprintln(os.Stderr, "Usage: protoc --api_linter_out=cfg_file=my_cfg_file,out_file=my_lint_output_file:. my_proto_file")
		os.Exit(1)
	}

	Run(opts, func(gen *Plugin) error {
		rt, err := createLinterRuntime(*cfgFile)
		if err != nil {
			return err
		}
		problems, err := lintFiles(rt, gen.Files())
		if err != nil {
			return err
		}
		var w io.Writer = os.Stderr
		if *outFile != "" {
			w = gen.NewGeneratedFile(*outFile)
		}
		var format formatFunc
		switch *outFormat {
		case "yaml":
			format = yaml.Marshal
		default:
			format = json.Marshal
		}
		return writeProblems(w, problems, format)
	})
}

type formatFunc func(interface{}) ([]byte, error)

func writeProblems(w io.Writer, problems []lint.Problem, f formatFunc) error {
	c, err := f(problems)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "%s\n", c)
	return nil
}

func createLinterRuntime(cfgFile string) (*lint.Runtime, error) {
	defaultCfg := lint.RuntimeConfig{
		IncludedPaths: []string{"**/*.proto"},
		RuleConfigs: map[string]lint.RuleConfig{
			"core": {
				Status:   lint.Enabled,
				Category: lint.Warning,
			},
		},
	}
	rt := lint.NewRuntime(defaultCfg)
	if cfgFile != "" {
		cfg, err := readRuntimeConfigs(cfgFile)
		if err != nil {
			return nil, err
		}
		rt.AddConfigs(cfg...)
	}
	rt.AddRules(rules.Rules().All()...)
	return rt, nil
}

func lintFiles(rt *lint.Runtime, files []*descriptorpb.FileDescriptorProto) ([]lint.Problem, error) {
	var problems []lint.Problem
	for _, proto := range files {
		req, err := lint.NewProtoRequest(proto)
		if err != nil {
			return nil, err
		}
		resp, err := rt.Run(req)
		if err != nil {
			return nil, err
		}
		for _, prob := range resp.Problems {
			problems = append(problems, prob)
		}
	}
	return problems, nil
}

func readRuntimeConfigs(cfgFile string) (lint.RuntimeConfigs, error) {
	f, err := os.Open(cfgFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	cfg, err := lint.ReadConfigsJSON(f)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
