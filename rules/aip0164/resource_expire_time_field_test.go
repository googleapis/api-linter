package aip0164

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestResourceExpireTimeField(t *testing.T) {
	for _, test := range []struct {
		desc          string
		MethodName    string
		ResourceField string
		problems      testutils.Problems
	}{
		{"Valid", `UndeleteBook`, `google.protobuf.Timestamp purge_time = 1;`, nil},
		{"Valid", `UndeleteBookLegacyExpireTime`, `google.protobuf.Timestamp expire_time = 1;`, nil},
		{"Invalid", `UndeleteBook`, ``, testutils.Problems{{Message: "Resources supporting soft delete"}}},
		{"IrrelevantNoSoftDelete", `GetBook`, ``, nil},
	} {
		t.Run(test.desc, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
				import "google/protobuf/timestamp.proto";
				service Library {
					rpc {{.MethodName}}({{.MethodName}}Request) returns (Book);
				}
				message Book {
					{{.ResourceField}}
				}
				message {{.MethodName}}Request {}
			`, test)
			message := file.GetMessageTypes()[0]
			problems := resourceExpireTimeField.Lint(file)
			if diff := test.problems.SetDescriptor(message).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}
