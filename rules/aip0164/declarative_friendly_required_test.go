package aip0164

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestDeclarativeFriendlyRequired(t *testing.T) {
	for _, test := range []struct {
		testName   string
		Style      string
		MethodName string
		problems   testutils.Problems
	}{
		{"ValidDF", `style: DECLARATIVE_FRIENDLY`, `UndeleteBook`, nil},
		{"ValidNotDF", ``, `GetBook`, nil},
		{"InvalidDF", `style: DECLARATIVE_FRIENDLY`, `GetBook`, testutils.Problems{{Message: "Book should have an Undelete"}}},
	} {
		t.Run(test.testName, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
				import "google/api/resource.proto";
				service Library {
					rpc {{.MethodName}}({{.MethodName}}Request) returns (Book);
				}
				message Book {
					option (google.api.resource) = {
						{{.Style}}
					};
				}
				message {{.MethodName}}Request {}
			`, test)
			message := file.GetMessageTypes()[0]
			problems := declarativeFriendlyRequired.Lint(file)
			if diff := test.problems.SetDescriptor(message).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}
