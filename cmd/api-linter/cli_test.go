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
			inputArgs: strings.Split("-config=config -out_format=json -out_path=out -proto_desc=proto_desc -proto_path=proto_path_a -proto_path=proto_path_b a.proto b.proto", " "),
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
