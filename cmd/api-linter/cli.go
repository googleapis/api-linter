package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/googleapis/api-linter/lint"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

func run(rules []lint.Rule, configs lint.RuntimeConfigs, args []string) error {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "fmt",
			Value: "yaml",
			Usage: "output format",
		},
		cli.StringFlag{
			Name:  "cfg",
			Value: "",
			Usage: "configuration file",
		},
		cli.StringFlag{
			Name:  "import_path",
			Value: ".",
			Usage: "protoc import path",
		},
		cli.StringFlag{
			Name:  "out",
			Value: "",
			Usage: "output file",
		},
	}
	app.Action = func(c *cli.Context) error {
		filenames := c.Args()
		if len(filenames) == 0 {
			return nil
		}

		p := protocParser{
			importPath: c.String("import_path"),
		}
		files, err := p.ParseProto(filenames...)
		if err != nil {
			return err
		}

		userConfigs, err := readConfigs(c.String("cfg"))
		if err != nil {
			return err
		}
		configs = append(configs, userConfigs...)
		l := newLinter(rules, configs)
		problems, err := l.LintProto(files)
		if err != nil {
			return err
		}

		w := os.Stdout
		if c.String("out") != "" {
			var err error
			w, err = os.Create(c.String("out"))
			if err != nil {
				return err
			}
			defer w.Close()
		}

		marshal := yaml.Marshal
		switch c.String("fmt") {
		case "json":
			marshal = json.Marshal
		}
		b, err := marshal(problems)
		if err != nil {
			return err
		}

		if _, err = w.Write(b); err != nil {
			return err
		}
		return nil
	}

	return app.Run(args)
}

func readConfigs(path string) (lint.RuntimeConfigs, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("readConfig: %s", err.Error())
	}
	defer f.Close()

	parse := lint.ReadConfigsJSON
	switch filepath.Ext(path) {
	case ".yaml":
		parse = lint.ReadConfigsYAML
	}
	return parse(f)
}
