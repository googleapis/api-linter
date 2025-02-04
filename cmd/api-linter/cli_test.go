package main

import (
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
			name: "AllFlags",
			inputArgs: []string{
				"--config=config",
				"--output-format=json",
				"-o=out",
				"--descriptor-set-in=proto_desc1",
				"--descriptor-set-in=proto_desc2",
				"--proto-path=proto_path_a",
				"-I=proto_path_b",
				"a.proto",
				"b.proto",
			},
			wantCli: &cli{
				ConfigPath:       "config",
				OutputPath:       "out",
				FormatType:       "json",
				ProtoDescPath:    []string{"proto_desc1", "proto_desc2"},
				ProtoImportPaths: []string{"proto_path_a", "proto_path_b"},
				ProtoFiles:       []string{"a.proto", "b.proto"},
			},
		},
		{
			name: "ExitStatusOnLintFailure",
			inputArgs: []string{
				"--set-exit-status",
			},
			wantCli: &cli{
				ExitStatusOnLintFailure: true,
				ProtoFiles:              []string{},
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
