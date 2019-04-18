package rules

import (
	"github.com/jgeewax/api-linter/lint"
	"testing"
)

func TestFieldNamesRule_ConformingFieldNames(t *testing.T) {
	pd, err := protoDescriptorProtoFromSource(`syntax = "proto2";

package google.apis.tools.analyzer.testprotos;

message Foo {
  optional string first_field_name = 1;

  optional string another_field_name = 2;
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
