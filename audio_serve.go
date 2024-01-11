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
	"io"
	"net/http"
	"os/exec"
)

type source string

func (s source) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Set some headers
	w.Header().Add("Content-Type", "audio/mpeg")
	w.Header().Add("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Add("Pragma", "no-cache")
	w.Header().Add("Expires", "0")
	w.Header().Add("Accept-Ranges", "none")
	w.Header().Add("Access-Control-Allow-Origin", "*")

	parecCMD := exec.Command("parec", "-d", string(s))
	ffmpegCMD := exec.Command(
		"ffmpeg",
		"-f", "s16le",
		"-ac", "2",
		"-i", "-",
		"-b:a", fmt.Sprintf("%dk", MP3BITRATE),
		"-f", "mp3", "-",
	)
	parecReader, parecWriter := io.Pipe()
	parecCMD.Stdout = parecWriter
	ffmpegCMD.Stdin = parecReader

	ffmpegReader, ffmpegWriter := io.Pipe()
	ffmpegCMD.Stdout = ffmpegWriter

	parecCMD.Start()
	ffmpegCMD.Start()

	io.Copy(w, ffmpegReader)

	if parecCMD.Process != nil {
		parecCMD.Process.Kill()
	}
	if ffmpegCMD.Process != nil {
		ffmpegCMD.Process.Kill()
	}
}
