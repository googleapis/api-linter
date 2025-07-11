package main

import (
	"os"
	"path/filepath"
	"sort"
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

func TestResolveImports(t *testing.T) {
	// Save the original working directory and restore it at the end of the test.
	originalCWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get original working directory: %v", err)
	}
	defer func() {
		if err := os.Chdir(originalCWD); err != nil {
			t.Errorf("Failed to restore original working directory: %v", err)
		}
	}()

	defaultSetup := func(t *testing.T, cwd, externalDir string) {}

	tests := []struct {
		name             string
		protoImportPaths func(cwd, externalDir string) []string
		setup            func(t *testing.T, cwd, externalDir string)
		want             func(externalDir string) []string
	}{
		{
			name: "NoProtoImportPaths",
			protoImportPaths: func(_, _ string) []string {
				return []string{}
			},
			setup: defaultSetup,
			want: func(_ string) []string {
				return []string{"."}
			},
		},
		{
			name: "ExplicitDot",
			protoImportPaths: func(_, _ string) []string {
				return []string{"."}
			},
			setup: defaultSetup,
			want: func(_ string) []string {
				return []string{"."}
			},
		},
		{
			name: "SubdirectoryOfCWD",
			protoImportPaths: func(_, _ string) []string {
				return []string{"test_dir"}
			},
			setup: func(t *testing.T, cwd, _ string) {
				if err := os.Mkdir(filepath.Join(cwd, "test_dir"), 0755); err != nil {
					t.Fatalf("Failed to create temp subdirectory: %v", err)
				}
			},
			want: func(_ string) []string {
				return []string{"."} // "test_dir" should be covered by "."
			},
		},
		{
			name: "SubdirectoryOfCWDWithDot",
			protoImportPaths: func(_, _ string) []string {
				return []string{".", "test_dir"}
			},
			setup: func(t *testing.T, cwd, _ string) {
				if err := os.Mkdir(filepath.Join(cwd, "test_dir"), 0755); err != nil {
					t.Fatalf("Failed to create temp subdirectory: %v", err)
				}
			},
			want: func(_ string) []string {
				return []string{"."}
			},
		},
		{
			name: "ExternalAbsolutePath",
			protoImportPaths: func(_, externalDir string) []string {
				return []string{externalDir}
			},
			setup: defaultSetup,
			want: func(externalDir string) []string {
				return []string{".", externalDir}
			},
		},
		{
			name: "MixedPaths",
			protoImportPaths: func(_, externalDir string) []string {
				return []string{"./relative_dir", externalDir, "test_dir"}
			},
			setup: func(t *testing.T, cwd, _ string) {
				if err := os.Mkdir(filepath.Join(cwd, "relative_dir"), 0755); err != nil {
					t.Fatalf("Failed to create relative_dir: %v", err)
				}
				if err := os.Mkdir(filepath.Join(cwd, "test_dir"), 0755); err != nil {
					t.Fatalf("Failed to create test_dir: %v", err)
				}
			},
			want: func(externalDir string) []string {
				return []string{".", externalDir}
			},
		},
		{
			name: "NonExistentRelativePath",
			protoImportPaths: func(_, _ string) []string {
				return []string{"non_existent_dir"}
			},
			setup: defaultSetup,
			want: func(_ string) []string {
				return []string{"."} // Should still resolve to just "." as non_existent_dir is relative to CWD
			},
		},
		{
			name: "CurrentDirAsAbsolutePath",
			protoImportPaths: func(cwd, _ string) []string {
				return []string{cwd}
			},
			setup: defaultSetup,
			want: func(_ string) []string {
				return []string{"."}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cwd := t.TempDir()
			externalDir := t.TempDir()

			if err := os.Chdir(cwd); err != nil {
				t.Fatalf("Failed to change directory to %q: %v", cwd, err)
			}

			test.setup(t, cwd, externalDir)

			protoImportPaths := test.protoImportPaths(cwd, externalDir)
			want := test.want(externalDir)

			got := resolveImports(protoImportPaths)
			sort.Strings(got)
			sort.Strings(want)

			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("resolveImports() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
