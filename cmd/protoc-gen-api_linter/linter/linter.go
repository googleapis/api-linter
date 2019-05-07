// Package linter provides an API linter that
// can be run as a plugin in protoc.
package linter

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	descriptorpb "github.com/golang/protobuf/v2/types/descriptor"
	"github.com/googleapis/api-linter/cmd/protoc-gen-api_linter/protogen"
	"github.com/googleapis/api-linter/lint"
	"gopkg.in/yaml.v2"
)

// Linter is a protoc plugin that generates linting results for
// API protobuf files.
//
// See github.com/googleapis/api-linter/protoc-gen-api_linter/main.go
// for how to use it.
type Linter struct {
	params map[string]string
	rt     *lint.Runtime
}

// New creates and returns a Linter with the given rules and configs.
func New(rules []lint.Rule, configs lint.RuntimeConfigs) *Linter {
	l := &Linter{
		params: make(map[string]string),
		rt:     lint.NewRuntime(),
	}
	l.rt.AddRules(rules...)
	l.rt.AddConfigs(configs...)
	return l
}

func (linter *Linter) hasParam(name string) bool {
	v, found := linter.params[name]
	return found && v != ""
}

func (linter *Linter) getParam(name string) string {
	if linter.hasParam(name) {
		return linter.params[name]
	}
	return ""
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
// and write them to an output file, if specified, or os.Stderr.
func (linter *Linter) Generate(gen *protogen.Plugin) error {
	if linter.hasParam("cfg_file") {
		configs, err := readConfigs(linter.getParam("cfg_file"))
		if err != nil {
			return err
		}
		linter.rt.AddConfigs(configs...)
	}

	problems, err := checkFiles(linter.rt, gen.Files())
	if err != nil {
		return err
	}
	var w io.Writer = os.Stderr
	if linter.hasParam("out_file") {
		w = gen.NewGeneratedFile(linter.getParam("out_file"))
	}
	format := linter.getParam("out_fmt")
	return writeProblems(w, problems, format)
}

func writeProblems(w io.Writer, problems []lint.Problem, format string) error {
	var f func(interface{}) ([]byte, error)
	switch format {
	case "yaml":
		f = yaml.Marshal
	default:
		f = json.Marshal
	}
	c, err := f(problems)
	if err != nil {
		return err
	}
	_, err = w.Write(c)
	return err
}

func checkFiles(rt *lint.Runtime, files []*descriptorpb.FileDescriptorProto) ([]lint.Problem, error) {
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

func readConfigs(file string) (lint.RuntimeConfigs, error) {
	f, err := os.Open(file)
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
