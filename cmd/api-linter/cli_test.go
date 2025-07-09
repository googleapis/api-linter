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

	tests := []struct {
		name             string
		protoImportPaths []string
		setupTempDir     func(t *testing.T) string // Function to set up a temporary directory
		want             []string
	}{
		{
			name:             "NoProtoImportPaths",
			protoImportPaths: []string{},
			want:             []string{"."},
		},
		{
			name:             "ExplicitDot",
			protoImportPaths: []string{"."},
			want:             []string{"."},
		},
		{
			name:             "SubdirectoryOfCWD",
			protoImportPaths: []string{"test_dir"},
			setupTempDir: func(t *testing.T) string {
				tmpDir := t.TempDir()
				if err := os.Mkdir(filepath.Join(tmpDir, "test_dir"), 0755); err != nil {
					t.Fatalf("Failed to create temp subdirectory: %v", err)
				}
				return tmpDir
			},
			want: []string{"."}, // "test_dir" should be covered by "."
		},
		{
			name:             "SubdirectoryOfCWDWithDot",
			protoImportPaths: []string{".", "test_dir"},
			setupTempDir: func(t *testing.T) string {
				tmpDir := t.TempDir()
				if err := os.Mkdir(filepath.Join(tmpDir, "test_dir"), 0755); err != nil {
					t.Fatalf("Failed to create temp subdirectory: %v", err)
				}
				return tmpDir
			},
			want: []string{"."},
		},
		{
			name:             "ExternalAbsolutePath",
			protoImportPaths: []string{"/tmp/external_proto"},
			setupTempDir: func(t *testing.T) string {
				tmpDir := t.TempDir()
				// Create a dummy external directory to ensure filepath.Abs works
				if err := os.MkdirAll("/tmp/external_proto", 0755); err != nil {
					t.Fatalf("Failed to create external temp directory: %v", err)
				}
				t.Cleanup(func() { os.RemoveAll("/tmp/external_proto") })
				return tmpDir
			},
			want: []string{".", "/tmp/external_proto"},
		},
		{
			name:             "MixedPaths",
			protoImportPaths: []string{"./relative_dir", "/tmp/external_dir", "test_dir"},
			setupTempDir: func(t *testing.T) string {
				tmpDir := t.TempDir()
				if err := os.Mkdir(filepath.Join(tmpDir, "relative_dir"), 0755); err != nil {
					t.Fatalf("Failed to create relative_dir: %v", err)
				}
				if err := os.Mkdir(filepath.Join(tmpDir, "test_dir"), 0755); err != nil {
					t.Fatalf("Failed to create test_dir: %v", err)
				}
				if err := os.MkdirAll("/tmp/external_dir", 0755); err != nil {
					t.Fatalf("Failed to create external_dir: %v", err)
				}
				t.Cleanup(func() { os.RemoveAll("/tmp/external_dir") })
				return tmpDir
			},
			want: []string{".", "/tmp/external_dir"},
		},
		{
			name:             "NonExistentRelativePath",
			protoImportPaths: []string{"non_existent_dir"},
			setupTempDir: func(t *testing.T) string {
				return t.TempDir() // No special setup needed, just a temp dir
			},
			want: []string{"."}, // Should still resolve to just "." as non_existent_dir is relative to CWD
		},
		{
			name:             "CurrentDirAsAbsolutePath",
			protoImportPaths: []string{}, // Will be set dynamically inside t.Run
			setupTempDir: func(t *testing.T) string {
				return t.TempDir() // No special setup needed, just a temp dir
			},
			want: []string{"."},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Set up a temporary directory for the test if a setup function is provided.
			// Change to this directory.
			var currentTestDir string
			if test.setupTempDir != nil {
				currentTestDir = test.setupTempDir(t)
			} else {
				currentTestDir = t.TempDir()
			}

			// Dynamically set protoImportPaths for this specific test case
			var actualProtoImportPaths []string
			if test.name == "CurrentDirAsAbsolutePath" {
				actualProtoImportPaths = []string{currentTestDir}
			} else {
				actualProtoImportPaths = test.protoImportPaths
			}

			if err := os.Chdir(currentTestDir); err != nil {
				t.Fatalf("Failed to change directory to %q: %v", currentTestDir, err)
			}

			got := resolveImports(actualProtoImportPaths)
			sort.Strings(got)
			sort.Strings(test.want)

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("resolveImports() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
