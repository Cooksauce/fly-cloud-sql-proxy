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
	userStr := os.Getenv("RUN_USER")
	sockPath := os.Getenv("RUN_GRANT_SOCK")
	fileInfo, err := os.Stat(sockPath)
	if err != nil {
		os.Stderr.WriteString(fmt.Errorf("failed to get socket info at %s: %w", sockPath, err).Error() + "\n")
		os.Exit(1)
		return
	}

	user, err := strconv.ParseInt(userStr, 10, 0)
	if err != nil {
		os.Stderr.WriteString(fmt.Errorf("failed to parse user %s: %w", userStr, err).Error() + "\n")
		os.Exit(1)
		return
	}

	permAll, _ := strconv.ParseUint("0777", 8, 32)

	mode := fileInfo.Mode()
	modeNew := mode | fs.FileMode(permAll)

	if err := os.Chmod(sockPath, modeNew); err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}
	os.Stdout.WriteString(fmt.Sprintf("changed %s mode from %s to %s", sockPath, mode, modeNew) + "\n")

	if err := syscall.Setuid(int(user)); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}

	cmd.Execute()
}

// func mainn() {
// 	optStr := os.Getenv("SOCKET_FWD_OPTS")
// 	if optStr == "" {
// 		cmd.Execute()
// 		return
// 	}

// 	var opts SocketForwardOptions
// 	if err := json.Unmarshal([]byte(optStr), &opts); err != nil {
// 		os.Stderr.WriteString(fmt.Errorf("failed to parse socket opts: %w", err).Error())
// 		os.Exit(1)
// 	}

// 	if opts.Root == "" || opts.User == "" || opts.Su == 0 {
// 		os.Stderr.WriteString(fmt.Errorf("invalid opts: root, user, and su must all be present").Error())
// 		os.Stderr.WriteString(fmt.Sprintf("%+v", opts))
// 		os.Exit(1)
// 	}

// 	if os.Getenv("SOCKET_FWD_CHILD") != "" {
// 		if err := syscall.Setuid(opts.Su); err != nil {
// 			os.Stderr.WriteString(err.Error())
// 			os.Exit(1)
// 		}
// 		cmd.Execute()
// 		return
// 	}

// 	l, err := flyBind()
// 	if err != nil {
// 		os.Stderr.WriteString(err.Error())
// 		os.Exit(1)
// 	}

// 	if l != nil {

// 	}

// 	cmd.Execute()
// }

// func execAsUser(opts *SocketForwardOptions) {
// 	env := os.Environ()

// 	cmd := exec.CommandContext(context.Background(), "/self", os.Args...)
// 	cmd.Env = []string{""}
// }

// func flyBind() (net.Listener, error) {
// 	// { root: /.fly/api, user: /.fwd/fly/api, su: 65532 }
// 	optStr := os.Getenv("SOCKET_FWD_OPTS")
// 	if optStr != "" {
// 		os.Stdout.WriteString("No SOCKET_FWD_OPTS found. Skipping")
// 		return nil, nil
// 	}

// 	var opts SocketForwardOptions
// 	if err := json.Unmarshal([]byte(optStr), &opts); err != nil {
// 		return nil, fmt.Errorf("failed to parse socket options: %s", optStr)
// 	}

// 	userListen, err := net.Listen("unix", opts.User)
// 	if err != nil {
// 		return nil, err
// 	}

// 	go func() {
// 		for {
// 			user, err := userListen.Accept()
// 			if err != nil {
// 				os.Stderr.WriteString(fmt.Sprintf("failed to accept user conn: %s", err.Error()))
// 			}
// 			go handleConnect(user, &opts)
// 		}
// 	}()

// 	return userListen, nil
// }

// func handleConnect(user net.Conn, opts *SocketForwardOptions) error {
// 	root, err := net.Dial("unix", opts.Root)
// 	if err != nil {
// 		return err
// 	}
// 	go doCopy(user, root)
// 	go doCopy(root, user)
// 	return nil
// }

// func doCopy(src net.Conn, dst net.Conn) {
// 	defer src.Close()
// 	defer dst.Close()
// 	io.Copy(dst, src)
// }

// type SocketForwardOptions struct {
// 	Root string `json:"root"`
// 	User string `json:"user"`
// 	Su   int    `json:"su"`
// }
