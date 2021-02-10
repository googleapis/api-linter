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
				import "google/protobuf/empty.proto";
				service Library {
					rpc {{.MethodName}}({{.MethodName}}Request) returns (Book);
					rpc DeleteBook(DeleteBookRequest) returns (google.protobuf.Empty);
				}
				message Book {
					option (google.api.resource) = {
						{{.Style}}
					};
				}
				message {{.MethodName}}Request {}
				message DeleteBookRequest {}
			`, test)
			message := file.GetMessageTypes()[0]
			problems := declarativeFriendlyRequired.Lint(file)
			if diff := test.problems.SetDescriptor(message).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}

	// Also test that undelete is not required if delete is not present.
	t.Run("ValidNoDelete", func(t *testing.T) {
		file := testutils.ParseProto3String(t, `
			import "google/api/resource.proto";
			service Library {
				rpc GetBook(GetBookRequest) returns (Book);
			}
			message Book {
				option (google.api.resource) = {
					style: DECLARATIVE_FRIENDLY
				};
			}
			message GetBookRequest {}
		`)
		problems := declarativeFriendlyRequired.Lint(file)
		if problems != nil && len(problems) > 0 {
			t.Errorf("Got %v, expected no problems.", problems)
		}
	})
}
