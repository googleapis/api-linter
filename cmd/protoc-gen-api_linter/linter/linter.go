package linter

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	descriptorpb "github.com/golang/protobuf/v2/types/descriptor"
	"github.com/googleapis/api-linter/cmd/protoc-gen-api_linter/protogen"
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules"
	"gopkg.in/yaml.v2"
)

// Linter is a protoc plugin that generates linting results for
// API protobuf files.
type Linter struct {
	params map[string]string
}

// New creates and returns a Linter.
func New() *Linter {
	return &Linter{
		params: make(map[string]string),
	}
}

func (linter *Linter) setParam(name, value string) error {
	if _, found := linter.params[name]; found {
		return fmt.Errorf("parameter %q is duplicated", name)
	}
	linter.params[name] = value
	return nil
}

// Options returns Options when creating a Plugin.
func (linter *Linter) Options() *protogen.Options {
	return &protogen.Options{
		ParamFunc: linter.setParam,
	}
}

// Generate generates linting results for the files in the Plugin
// and write them to an output file, if exists, or os.Stderr.
func (linter *Linter) Generate(gen *protogen.Plugin) error {
	rt, err := createLinterRuntime(linter.params["cfg_file"])
	if err != nil {
		return err
	}
	problems, err := lintFiles(rt, gen.Files())
	if err != nil {
		return err
	}
	var w io.Writer = os.Stderr
	if linter.params["out_file"] != "" {
		w = gen.NewGeneratedFile(linter.params["out_file"])
	}
	var format formatFunc
	switch linter.params["out_format"] {
	case "yaml":
		format = yaml.Marshal
	default:
		format = json.Marshal
	}
	return writeProblems(w, problems, format)
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
