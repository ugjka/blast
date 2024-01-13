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

	"github.com/davecgh/go-spew/spew"
)

type source string

func (s source) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if *headers {
		spew.Dump(r.Method)
		spew.Dump(r.Header)
	}

	// Set some headers
	w.Header().Add("Cache-Control", "no-cache, no-store")
	w.Header().Add("Pragma", "no-cache")
	w.Header().Add("Expires", "0")
	w.Header().Add("Accept-Ranges", "none")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	// handle devices like Samsung TVs
	if r.Header.Get("GetContentFeatures.DLNA.ORG") == "1" {
		f := dlnaContentFeatures{
			profileName:     "MP3",
			supportTimeSeek: false,
			supportRange:    false,
			flags:           DLNA_ORG_FLAG_STREAMING_TRANSFER_MODE | DLNA_ORG_FLAG_DLNA_V15,
		}
		w.Header().Set("ContentFeatures.DLNA.ORG", f.String())
	}

	const yearSeconds = 365 * 24 * 60 * 60
	const yearBytes = yearSeconds * (MP3BITRATE / 8) * 1000
	if r.Header.Get("Getmediainfo.sec") == "1" {
		w.Header().Set("MediaInfo.sec", fmt.Sprintf("SEC_Duration=%d", yearSeconds*1000))
	}
	w.Header().Add("Content-Length", fmt.Sprint(yearBytes))
	w.Header().Add("Content-Type", "audio/mpeg")

	if r.Method == http.MethodHead {
		w.WriteHeader(http.StatusOK)
		return
	}

	parecCMD := exec.Command("parec", "-d", string(s), "-n", "blast-rec")
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
