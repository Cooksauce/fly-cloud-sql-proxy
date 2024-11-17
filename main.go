// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build !windows
// +build !windows

package main

import (
	"fmt"
	"io/fs"
	"os"
	"strconv"
	"syscall"

	"github.com/GoogleCloudPlatform/cloud-sql-proxy/v2/cmd"
)

func main() {
	if err := grantSocket(os.Getenv("GRANT_SOCK")); err != nil {
		os.Stderr.WriteString(fmt.Errorf("failed to grant socket: %w", err).Error() + "\n")
		os.Exit(1)
	}

	if err := changeUser(os.Getenv("SWITCH_USER")); err != nil {
		os.Stderr.WriteString(fmt.Errorf("failed to switch user: %w", err).Error() + "\n")
		os.Exit(1)
	}

	cmd.Execute()
}

func changeUser(userStr string) error {
	user, err := strconv.ParseInt(userStr, 10, 0)
	if err != nil {
		return fmt.Errorf("failed to parse user %s: %w", userStr, err)
	}

	if err := syscall.Setuid(int(user)); err != nil {
		return fmt.Errorf("failed to set user %s: %w", userStr, err)
	}

	return nil
}

func grantSocket(sockPath string) error {
	fileInfo, err := os.Stat(sockPath)
	if err != nil {
		return fmt.Errorf("failed to get socket info at %s: %w", sockPath, err)
	}

	permAll, _ := strconv.ParseUint("0777", 8, 32)

	mode := fileInfo.Mode()
	modeNew := mode | fs.FileMode(permAll)

	if err := os.Chmod(sockPath, modeNew); err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}
	os.Stdout.WriteString(fmt.Sprintf("changed %s mode from %s to %s", sockPath, mode, modeNew) + "\n")
	return nil
}
