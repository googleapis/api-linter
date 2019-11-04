package main

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewCli(t *testing.T) {
	tests := []struct {
		name      string
		inputArgs []string
		wantCli   *cli
	}{
		{
			name:      "AllFlags",
			inputArgs: strings.Split("--config=config --output-format=json -o=out --proto-descriptor-set=proto_desc -I=proto_path_a -I=proto_path_b a.proto b.proto", " "),
			wantCli: &cli{
				ConfigPath:       "config",
				OutputPath:       "out",
				FormatType:       "json",
				ProtoDescPath:    "proto_desc",
				ProtoImportPaths: []string{"proto_path_a", "proto_path_b", "."},
				ProtoFiles:       []string{"a.proto", "b.proto"},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotCli := newCli(test.inputArgs)
			if diff := cmp.Diff(gotCli, test.wantCli); diff != "" {
				t.Errorf("newCli() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
