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
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"slices"

	"github.com/davecgh/go-spew/spew"
)

type stream struct {
	sink         string
	mime         string
	format       string
	bitrate      int
	chunk        int
	printheaders bool
	contentfeat  dlnaContentFeatures
}

func (s stream) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.printheaders {
		spew.Fdump(os.Stderr, r.Proto)
		spew.Fdump(os.Stderr, r.RemoteAddr)
		spew.Fdump(os.Stderr, r.URL)
		spew.Fdump(os.Stderr, r.Method)
		spew.Fdump(os.Stderr, r.Header)
	}
	// Set some headers
	w.Header().Add("Cache-Control", "No-Cache, No-Store")
	w.Header().Add("Pragma", "No-Cache")
	w.Header().Add("Expires", "0")
	w.Header().Add("User-Agent", "Blast-DLNA UPnP/1.0 DLNADOC/1.50")
	// handle devices like Samsung TVs
	if r.Header.Get("GetContentFeatures.DLNA.ORG") == "1" {
		w.Header().Set("ContentFeatures.DLNA.ORG", s.contentfeat.String())
	}

	var yearSeconds = 365 * 24 * 60 * 60
	if r.Header.Get("Getmediainfo.sec") == "1" {
		w.Header().Set("MediaInfo.sec", fmt.Sprintf("SEC_Duration=%d", yearSeconds*1000))
	}
	w.Header().Add("Content-Type", s.mime)

	flusher, ok := w.(http.Flusher)
	chunked := ok && r.Proto == "HTTP/1.1"

	if !chunked {
		var yearBytes = yearSeconds * (s.bitrate / 8) * 1000
		w.Header().Add("Content-Length", fmt.Sprint(yearBytes))
	}

	if r.Method == http.MethodHead {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	parecCMD := exec.Command("parec", "-d", s.sink, "-n", "blast-rec")
	parecErrBuf := bytes.NewBuffer(nil)

	ffargs := []string{"-loglevel", "error",
		"-f", "s16le",
		"-ac", "2",
		"-i", "-",
		"-f", s.format, "-"}
	if s.bitrate != 0 {
		ffargs = slices.Insert(
			ffargs,
			len(ffargs)-3,
			"-b:a", fmt.Sprintf("%dk", s.bitrate),
		)
	}
	if s.format == "m4a" || s.format == "mp4" {
		ffargs = slices.Insert(
			ffargs,
			len(ffargs)-3,
			"-movflags", "+faststart",
		)
	}
	ffmpegCMD := exec.Command("ffmpeg", ffargs...)

	ffmpegErrBuf := bytes.NewBuffer(nil)
	ffmpegCMD.Stderr = ffmpegErrBuf

	parecReader, parecWriter := io.Pipe()
	parecCMD.Stdout = parecWriter
	ffmpegCMD.Stdin = parecReader

	ffmpegReader, ffmpegWriter := io.Pipe()
	ffmpegCMD.Stdout = ffmpegWriter

	err := parecCMD.Start()
	if err != nil {
		log.Printf("parec failed: %v", err)
		return
	}
	err = ffmpegCMD.Start()
	if err != nil {
		log.Printf("ffmpeg failed: %v", err)
		return
	}
	if chunked {
		var (
			err error
			n   int
		)
		buf := make([]byte, (s.bitrate/8)*1000*s.chunk)
		if s.bitrate == 0 {
			buf = make([]byte, 44100*16*2)
		}
		for {
			n, err = ffmpegReader.Read(buf)
			if err != nil {
				break
			}
			_, err = w.Write(buf[:n])
			if err != nil {
				break
			}
			flusher.Flush()
		}
	} else {
		io.Copy(w, ffmpegReader)
	}

	if parecCMD.Process != nil {
		parecCMD.Process.Kill()
	}
	if ffmpegCMD.Process != nil {
		ffmpegCMD.Process.Kill()
	}
	if parecErrBuf.Len() > 0 {
		log.Printf("parec stderr: %s", parecErrBuf)
	}
	if ffmpegErrBuf.Len() > 0 {
		log.Printf("ffmpeg stderr: %s", ffmpegErrBuf)
	}
}
