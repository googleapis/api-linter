// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// This tool reads from the `.git` directory and determines the upcoming
// version tag and changelog.
//
// Usage (see .github/workflows/release.yaml):
//   go run ./.github/release-tool
package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func main() {
	// Get the most recent release version.
	lastTaggedRev := mustExec("git", "rev-list", "--tags", "--max-count=1")
	lastVersion := versionFromString(mustExec("git", "describe", "--tags", lastTaggedRev))

	// Get the changelog between the most recent release version and now.
	cl := newChangelog(
		mustExec("git", "log", fmt.Sprintf("%s..HEAD", lastVersion), "--oneline", "--pretty=format:%s"),
	)

	// Dump output.
	nextVersion := cl.incrVersion(lastVersion)
	expr := "::set-output name=%s::%s\n"
	if lastVersion != nextVersion {
		fmt.Printf("New version: %s\n", nextVersion)
		fmt.Printf(expr, "version", nextVersion.String()[1:])
		fmt.Printf(expr, "release_notes", cl.notes())
	} else {
		fmt.Printf("No changes from %s to now.", lastVersion)
	}
}

func mustExec(cmd string, args ...string) string {
	out, err := exec.Command(cmd, args...).CombinedOutput()
	if err != nil {
		log.Fatalf("exec failed: %s\n%s", out, err)
	}
	return strings.TrimSpace(string(out))
}
