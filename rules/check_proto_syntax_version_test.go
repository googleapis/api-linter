package rules

import (
	"fmt"
	"testing"

	"github.com/googleapis/api-linter/rules/testdata"
)

func TestProtoVersionRule(t *testing.T) {
	tmpl := testdata.MustCreateTemplate(`syntax = "{{.Syntax}}";`)

	tests := []struct {
		Syntax     string
		numProblem int
		suggestion string
		startLine  int
	}{
		{"proto3", 0, "", 1},
		{"proto2", 1, "proto3", 1},
	}

	rule := checkProtoVersion()
	for _, test := range tests {
		req := testdata.MustCreateRequestFromTemplate(tmpl, test)

		errPrefix := fmt.Sprintf("Check syntax `%s`", test.Syntax)

		resp, err := rule.Lint(req)
		if err != nil {
			t.Errorf("%s: Lint return error %v", errPrefix, err)
		}

		if got, want := len(resp.Problems), test.numProblem; got != want {
			t.Errorf("%s: got %d problems, but want %d", errPrefix, got, want)
		}

		if len(resp.Problems) > 0 {
			if got, want := resp.Problems[0].Suggestion, test.suggestion; got != want {
				t.Errorf("%s: got suggestion '%s', but want '%s'", errPrefix, got, want)
			}
			if got, want := resp.Problems[0].Location.Start.Line, test.startLine; got != want {
				t.Errorf("%s: got location starting with %d, but want %d", errPrefix, got, want)
			}
		}
	}
}
