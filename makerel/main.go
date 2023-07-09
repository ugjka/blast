// MIT+NoAI License
//
// Copyright (c) 2023 ugjka <ugjka@proton.me>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights///
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//
// This code may not be used to train artificial intelligence computer models
// or retrieved by artificial intelligence software or hardware.
package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	targets := []struct {
		os   string
		arch []string
	}{
		{"linux",
			[]string{
				"386",
				"amd64",
				"arm",
				"arm64",
				"loong64",
				"mips",
				"mips64",
				"mips64le",
				"mipsle",
				"ppc64",
				"ppc64le",
				"riscv64",
				"s390x",
			},
		},
	}

	for _, t := range targets {
		for _, arch := range t.arch {
			build := exec.Command("go", "build")
			build.Stderr = os.Stderr
			build.Stdout = os.Stdout
			build.Env = append(os.Environ(), "GOOS="+t.os, "GOARCH="+arch)
			if err := build.Run(); err != nil {
				panic(err)
			}
			zip := exec.Command("zip", "")
			zip.Stderr = os.Stderr
			zip.Stdout = os.Stdout
			zip.Args = []string{"-1", fmt.Sprintf("blast_%s_%s.zip", t.os, arch), "blast", "LICENSE", "README.md"}
			if err := zip.Run(); err != nil {
				panic(err)
			}
			if err := os.Remove("blast"); err != nil {
				panic(err)
			}
		}
	}
}
