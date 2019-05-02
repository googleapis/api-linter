package corp

import (
	"github.com/googleapis/api-linter/rules/rulestest"
	"testing"
)

func TestDoNotUseAny(t *testing.T) {
	rule := doNotUseAny()

	tpl := rulestest.MustCreateTemplate(`
syntax = "proto3";

package foo;

{{range .Imports}} import "{{.}}"; {{end}}

message Foo {
	string bar = 1;
}
`)

	tests := []struct {
		Imports     []string
		numProblems int
	}{
		{[]string{anyPath}, 1},
		{nil, 0},
	}

	for ind, test := range tests {
		req := rulestest.MustCreateRequestFromTemplate(tpl, test, "testdata/test_deps.protoset")

		p, err := rule.Lint(req)

		if err != nil {
			t.Errorf("Test #%d: Lint() returned an error: %v", ind+1, err)
		}

		if len(p.Problems) != test.numProblems {
			t.Errorf(
				"Test #%d: Lint() returned %d problems; want %d. Problems: %+v",
				ind+1, len(p.Problems), test.numProblems, p.Problems,
			)
		}
	}
}
