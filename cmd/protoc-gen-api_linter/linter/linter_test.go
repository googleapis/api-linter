package linter_test

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/googleapis/api-linter/cmd/protoc-gen-api_linter/linter"
	"github.com/googleapis/api-linter/cmd/protoc-gen-api_linter/protogen"
)

func TestAPILinter_OutputToFile(t *testing.T) {
	workdir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(workdir)

	tempfile := filepath.Join(workdir, "lint_test_out.json")
	runLinter(t, workdir, "out_file="+filepath.Base(tempfile), "testdata/test.proto")

	f, err := os.Open(tempfile)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	content, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(content), "message") {
		t.Errorf("Linting result: %q does not contains linting results with 'message'", content)
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
		protogen.Run(linter.New())
		os.Exit(0)
	}
}
