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
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/huin/goupnp"
)

// open in firewall
const STREAMPORT = 9000

// lame encoder bitrate
const MP3BITRATE = 320

const BLASTMONITOR = "blast.monitor"

func main() {
	// check for dependencies
	exes := []string{
		"pactl",
		"parec",
		"ffmpeg",
	}
	for _, exe := range exes {
		if _, err := exec.LookPath(exe); err != nil {
			stderr(err)
		}
	}

	var blastSinkID []byte
	var isPlaying bool
	var DLNADevice *goupnp.MaybeRootDevice
	// trap ctrl+c and kill
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		fmt.Println()
		if blastSinkID != nil {
			log.Println("unloading the blast sink")
			exec.Command("pactl", "unload-module", string(blastSinkID)).Run()
		}
		if isPlaying {
			log.Println("stopping av1transport and exiting")
			AV1Stop(DLNADevice.Location)
		}
		fmt.Println("terminated...")
		os.Exit(0)
	}()

	DLNADevice = chooseUPNPDevice()
	fmt.Println("----------")

	audioSource := chooseAudioSource()
	// on-demand handling of blast sink
	if string(audioSource) == BLASTMONITOR {
		blastSink := exec.Command("pactl", "load-module", "module-null-sink", "sink_name=blast")
		var err error
		blastSinkID, err = blastSink.Output()
		stderr(err)
		blastSinkID = bytes.TrimSpace(blastSinkID)
	}

	fmt.Println("----------")
	streamAddress := chooseStreamIP()
	fmt.Println("----------")

	log.Printf("starting the stream on port %d (configure your firewall if necessary)", STREAMPORT)

	mux := http.NewServeMux()
	mux.Handle("/stream", audioSource)
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", STREAMPORT),
		ReadTimeout:  -1,
		WriteTimeout: -1,
		Handler:      mux,
	}
	go func() {
		err := httpServer.ListenAndServe()
		stderr(err)
	}()
	// detect when the stream server is up
	for {
		_, err := net.Dial("tcp", fmt.Sprintf(":%d", STREAMPORT))
		if err == nil {
			break
		}
	}

	var streamURL string
	if streamAddress.To4() != nil {
		streamURL = fmt.Sprintf("http://%s:%d/stream", streamAddress, STREAMPORT)
	} else {
		streamURL = fmt.Sprintf("http://[%s]:%d/stream", streamAddress, STREAMPORT)
	}
	log.Printf("stream URI: %s\n", streamURL)
	log.Println("setting av1transport URI and playing")
	AV1SetAndPlay(DLNADevice.Location, streamURL)
	isPlaying = true
	select {}
}
