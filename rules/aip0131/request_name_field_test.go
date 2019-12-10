package aip0131

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestRequestNameFieldType(t *testing.T) {
	tests := []struct {
		name          string
		MessageName   string
		NameFieldType string
		problems      testutils.Problems
	}{
		{"StringNameFieldType_Valid", "GetBookRequest", "string", nil},
		{"BytesNameFieldType_Invalid", "GetBookRequest", "bytes", testutils.Problems{{Suggestion: "string"}}},
		{"NotGetRequest_BytesNameFieldType_Valid", "SomeMessage", "bytes", nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
			message {{.MessageName}} {
				{{.NameFieldType}} name = 1;
			}`, test)

			problems := requestNameField.Lint(f)
			if diff := test.problems.SetDescriptor(f.GetMessageTypes()[0].GetFields()[0]).Diff(problems); diff != "" {
				t.Errorf("Problems did not match: %v", diff)
			}
		})
	}
}
