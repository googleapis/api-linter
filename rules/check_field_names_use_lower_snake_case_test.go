package rules

import (
	"github.com/jgeewax/api-linter/lint"
	"testing"
)

func TestFieldNamesRule_ConformingFieldNames(t *testing.T) {
	pd, err := protoDescriptorProtoFromSource(`syntax = "proto2";

package google.apis.tools.analyzer.testprotos;

message Outer {
  message Middle {
      optional string middle_field_name = 1;
  }

  optional string outer_field_name = 2;

  enum NestedEnum {
    FOO = 1;
  }
  oneof outer_oneof_field {
    string outer_oneof_field_name = 3;
  }

  extensions 100 to 199;
}`)

	if err != nil {
		t.Fatalf("Error generating proto descriptor: %v", err)
	}

	rules, err := lint.NewRules(checkNamingFormats())

	if err != nil {
		t.Errorf("Error returned when creating Rules: %v", err)
	}

	req, err := lint.NewProtoFileRequest(pd)

	if err != nil {
		t.Errorf("Error returned when creating ProtoFileRequest: %v", err)
	}

	resp, err := lint.Run(rules, req)

	if len(resp.Problems) > 0 {
		t.Errorf("Expecting no problems, got %d", len(resp.Problems))
	}
}
