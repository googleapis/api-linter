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
	"log"
	"strconv"
	"strings"
)

type version struct {
	Major int `json:"major"`
	Minor int `json:"minor"`
	Patch int `json:"patch"`
}

func versionFromString(ver string) *version {
	// ver will be "vM.M.P", always. We will just panic if not.
	split := strings.Split(ver[1:], ".")
	return &version{
		Major: toInt(split[0]),
		Minor: toInt(split[1]),
		Patch: toInt(split[2]),
	}
}

func (v version) incrMajor() *version {
	return &version{Major: v.Major + 1, Minor: 0, Patch: 0}
}

func (v version) incrMinor() *version {
	return &version{Major: v.Major, Minor: v.Minor + 1, Patch: 0}
}

func (v version) incrPatch() *version {
	return &version{Major: v.Major, Minor: v.Minor, Patch: v.Patch + 1}
}

func (v version) String() string {
	return fmt.Sprintf("v%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func toInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("%s", err)
	}
	return i
}
