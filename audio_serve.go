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
	"os"
	"os/exec"

	"github.com/davecgh/go-spew/spew"
)

type source string

func (s source) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if *headers {
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
		f := dlnaContentFeatures{
			profileName:     "MP3",
			supportTimeSeek: false,
			supportRange:    false,
			flags: DLNA_ORG_FLAG_DLNA_V15 |
				DLNA_ORG_FLAG_CONNECTION_STALL |
				DLNA_ORG_FLAG_STREAMING_TRANSFER_MODE |
				DLNA_ORG_FLAG_BACKGROUND_TRANSFERT_MODE,
		}
		w.Header().Set("ContentFeatures.DLNA.ORG", f.String())
	}

	var yearSeconds = 365 * 24 * 60 * 60
	if r.Header.Get("Getmediainfo.sec") == "1" {
		w.Header().Set("MediaInfo.sec", fmt.Sprintf("SEC_Duration=%d", yearSeconds*1000))
	}
	w.Header().Add("Content-Type", "audio/mpeg")

	flusher, ok := w.(http.Flusher)
	chunked := ok && r.Proto == "HTTP/1.1"

	if !chunked {
		var yearBytes = yearSeconds * (MP3BITRATE / 8) * 1000
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
	if chunked {
		var (
			err error
			n   int
		)
		buf := make([]byte, (MP3BITRATE/8)*1000*CHUNK_SECONDS)
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
}
