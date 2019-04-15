package lint_test

import (
	"testing"

	"github.com/jgeewax/api-linter/lint"
	"github.com/jgeewax/api-linter/lint/mocks"
)

func TestRunAll(t *testing.T) {
	r1 := new(mocks.Rule)
	r2 := new(mocks.Rule)
	req := lint.Request{}

	resp1 := lint.Response{
		Problems: []lint.Problem{
			lint.Problem{},
		},
	}

	resp2 := lint.Response{
		Problems: []lint.Problem{
			lint.Problem{},
			lint.Problem{},
		},
	}

	r1.On("Lint", req).Return(resp1, nil)
	r2.On("Lint", req).Return(resp2, nil)
	r1.On("Name").Return("a")
	r2.On("Name").Return("b")

	rules, _ := lint.NewRules(r1, r2)
	gotResp, err := lint.Run(*rules, req)
	if err != nil {
		t.Errorf("Run: returns error %v, but wanted none", err)
	}
	if len(gotResp.Problems) != 3 {
		t.Errorf("Run: returns %d problems, but wanted 3", len(gotResp.Problems))
	}
}
