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
		{
			name: "VersionFlag",
			inputArgs: []string{
				"--version",
			},
			wantCli: &cli{
				VersionFlag: true,
				ProtoFiles:  []string{},
			},
		},
		{
			name: "ListRulesFlag",
			inputArgs: []string{
				"--list-rules",
				"--output-format=json", // Should also capture output format
			},
			wantCli: &cli{
				ListRulesFlag: true,
				FormatType:    "json",
				ProtoFiles:    []string{},
			},
		},
		{
			name: "DebugFlag",
			inputArgs: []string{
				"--debug",
			},
			wantCli: &cli{
				DebugFlag:  true,
				ProtoFiles: []string{},
			},
		},
		{
			name: "IgnoreCommentDisablesFlag",
			inputArgs: []string{
				"--ignore-comment-disables",
			},
			wantCli: &cli{
				IgnoreCommentDisablesFlag: true,
				ProtoFiles:                []string{},
			},
		},
		{
			name: "EnableAndDisableRules",
			inputArgs: []string{
				"--enable-rule=aip-123",
				"--enable-rule=aip-456",
				"--disable-rule=aip-789",
			},
			wantCli: &cli{
				EnabledRules:  []string{"aip-123", "aip-456"},
				DisabledRules: []string{"aip-789"},
				ProtoFiles:    []string{},
			},
		},
		{
			name: "CombinedFlags",
			inputArgs: []string{
				"--config=my-config.yaml",
				"--output-format=yaml",
				"-o=results.yml",
				"--set-exit-status",
				"--proto-path=include",
				"--descriptor-set-in=desc.pb",
				"--enable-rule=aip-999",
				"--debug",
				"file1.proto",
				"file2.proto",
			},
			wantCli: &cli{
				ConfigPath:              "my-config.yaml",
				FormatType:              "yaml",
				OutputPath:              "results.yml",
				ExitStatusOnLintFailure: true,
				ProtoImportPaths:        []string{"include"},
				ProtoDescPath:           []string{"desc.pb"},
				EnabledRules:            []string{"aip-999"},
				DisabledRules:           nil,
				DebugFlag:               true,
				ProtoFiles:              []string{"file1.proto", "file2.proto"},
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
