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

package main

import (
	"fmt"
	"strings"
)

type commit struct {
	category string
	message  string
	breaking bool
}

func newCommit(line string) *commit {
	var cat, msg string
	breaking := false
	lineSplit := strings.SplitN(line, ":", 2)
	if len(lineSplit) > 1 {
		cat = strings.TrimSpace(lineSplit[0])
		if strings.HasSuffix(cat, "!") {
			breaking = true
			cat = strings.TrimSuffix(cat, "!")
		}
		if scope := strings.Index(cat, "("); scope != -1 {
			cat = cat[:scope]
		}
		msg = strings.TrimSpace(lineSplit[1])
	} else {
		cat = "unknown"
		msg = strings.TrimSpace(lineSplit[0])
	}
	return &commit{category: cat, message: msg, breaking: breaking}
}

func (c *commit) String() string {
	return c.message
}

type changelog struct {
	features       []*commit
	fixes          []*commit
	docs           []*commit
	otherVisible   []*commit
	otherInvisible []*commit
	breaking       bool
}

func newChangelog(gitlog string) *changelog {
	breaking := false
	var feat, fix, docs, vis, invis []*commit
	for _, line := range strings.Split(gitlog, "\n") {
		cmt := newCommit(line)
		if cmt.breaking {
			breaking = true
		}
		switch cmt.category {
		case "feat":
			feat = append(feat, cmt)
		case "fix":
			fix = append(fix, cmt)
		case "docs":
			docs = append(docs, cmt)
		case "refactor":
			vis = append(vis, cmt)
		case "chore":
			vis = append(vis, cmt)
		default:
			invis = append(invis, cmt)
		}
	}
	return &changelog{
		features:       feat,
		fixes:          fix,
		docs:           docs,
		otherVisible:   vis,
		otherInvisible: invis,
		breaking:       breaking,
	}
}

func (cl *changelog) incrVersion(v *version) *version {
	if cl.breaking {
		return v.incrMajor()
	}
	if len(cl.features) > 0 {
		return v.incrMinor()
	}
	if len(cl.fixes) > 0 || len(cl.otherVisible) > 0 {
		return v.incrPatch()
	}
	return v
}

func (cl *changelog) notes() string {
	section := func(title string, commits []*commit) string {
		if len(commits) > 0 {
			// %0A is newline: https://github.community/t/set-output-truncates-multiline-strings/16852
			answer := fmt.Sprintf("## %s%%0A%%0A", title)
			for _, cmt := range commits {
				answer += fmt.Sprintf("- %s%%0A", cmt)
			}
			return answer + "%0A"
		}
		return ""
	}
	return strings.TrimSuffix(
		section("Features", cl.features)+section("Fixes", cl.fixes)+section("Documentation", cl.docs)+section("Other", cl.otherVisible),
		"%0A",
	)
}
