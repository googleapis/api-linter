package linter_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/googleapis/api-linter/cmd/protoc-gen-api_linter/linter"
	"github.com/googleapis/api-linter/cmd/protoc-gen-api_linter/protogen"
	"github.com/googleapis/api-linter/lint"
	corerules "github.com/googleapis/api-linter/rules"
)

func TestAPILinter_WithConfigFile(t *testing.T) {
	workdir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(workdir)

	tests := []struct {
		cfgfile string
		want    string
	}{
		{
			// By default, category is "warning".
			"",
			`"category":"warning"`,
		},
		{
			// Category will be overrode to "error".
			"testdata/test_rule_configs.json",
			`"category":"error"`,
		},
	}
	for _, test := range tests {
		outfile := "test.out"
		params := fmt.Sprintf("out_file=%s,cfg_file=%s", outfile, test.cfgfile)
		runLinter(t, workdir, params, "testdata/test.proto")

		f, err := os.Open(filepath.Join(workdir, outfile))
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		content, err := ioutil.ReadAll(f)
		if err != nil {
			t.Fatal(err)
		}

		if !strings.Contains(string(content), test.want) {
			t.Errorf("Linting result: %q does not contain linting results with %q", content, test.want)
		}
	}
}

// runLinter invokes protoc on the linter plugin.
func runLinter(t *testing.T, workdir, params string, args ...string) {
	cmd := exec.Command("protoc", "--plugin=protoc-gen-api_linter="+os.Args[0])
	cmd.Args = append(cmd.Args, "--api_linter_out="+params+":"+workdir, "-Itestdata")
	cmd.Args = append(cmd.Args, args...)
	cmd.Env = append(os.Environ(), "RUN_AS_PROTOC_PLUGIN=1")
	out, err := cmd.CombinedOutput()
	t.Log("RUNNING: ", strings.Join(cmd.Args, " "))
	if len(out) > 0 {
		t.Log(string(out))
	}
	if err != nil {
		t.Fatalf("protoc: %v", err)
	}
}

func init() {
	if os.Getenv("RUN_AS_PROTOC_PLUGIN") != "" {
		defaultConfigs := lint.RuntimeConfigs{
			lint.RuntimeConfig{
				IncludedPaths: []string{"**/*.proto"},
				RuleConfigs: map[string]lint.RuleConfig{
					"core": {
						Status:   lint.Enabled,
						Category: lint.Warning,
					},
				},
			},
		}
		var rules []lint.Rule
		rules = append(rules, corerules.Rules().All()...)
		protogen.Run(linter.New(rules, defaultConfigs))
		os.Exit(0)
	}
}
