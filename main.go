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
	_ "embed"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/davecgh/go-spew/spew"
	"github.com/huin/goupnp"
	"github.com/huin/goupnp/dcps/av1"
)

const (
	BLASTMONITOR = "blast.monitor"
	LOGO_PATH    = "logo.png"
)

//go:embed logo.png
var logobytes []byte

var logblast = new(bool)

func main() {
	// check for dependencies
	exes := []string{
		"pactl",
		"parec",
		"ffmpeg",
	}
	for _, exe := range exes {
		if _, err := exec.LookPath(exe); err != nil {
			if err != nil {
				fmt.Fprintln(os.Stderr, "dependency:", err)
				os.Exit(1)
			}
		}
	}
	device := flag.String("device", "", "dlna device's friendly name")
	source := flag.String("source", "", "audio source (pactl list sources short | cut -f2)")
	ip := flag.String("ip", "", "host ip address")
	port := flag.Int("port", 9000, "stream port")
	chunk := flag.Int("chunk", 1, "chunk size in seconds")
	bitrate := flag.Int("bitrate", 320, "audio format bitrate")
	format := flag.String("format", "mp3", "stream audio format")
	mime := flag.String("mime", "audio/mpeg", "stream mime type")
	useaac := flag.Bool("useaac", false, "use aac audio")
	useflac := flag.Bool("useflac", false, "use flac audio")
	uselpcm := flag.Bool("uselpcm", false, "use lpcm audio")
	usewav := flag.Bool("usewav", false, "use wav audio")
	bits := flag.Int("bits", 16, "audio bitdepth")
	rate := flag.Int("rate", 44100, "audio sample rate")
	channels := flag.Int("channels", 2, "audio channels")
	dummy := flag.Bool("dummy", false, "only serve content")
	debug := flag.Bool("debug", false, "print debug info")
	headers := flag.Bool("headers", false, "print request headers")
	logblast = flag.Bool("log", false, "log parec and ffmpeg stderr")
	nochunked := flag.Bool("nochunked", false, "disable chunked tranfer endcoding")
	bige := flag.Bool("bige", false, "use big endian for capture and lpcm format")

	flag.Parse()

	var (
		blastSinkID []byte
		isPlaying   bool
		DLNADevice  *goupnp.MaybeRootDevice
		err         error
	)

	// trap ctrl+c and kill and terminal hang up
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	cleanup := func() {
		if blastSinkID != nil {
			log.Println("unloading the blast sink")
			exec.Command("pactl", "unload-module", string(blastSinkID)).Run()
		}
	}

	go func() {
		<-sig
		fmt.Println()
		cleanup()
		if isPlaying && !*dummy {
			log.Println("stopping av1transport and exiting")
			AV1Stop(DLNADevice.Location)
		}
		fmt.Println("terminated...")
		os.Exit(0)
	}()
	if !*dummy {
		DLNADevice, err = chooseUPNPDevice(*device)
		if err != nil {
			fmt.Fprintln(os.Stderr, "upnp:", err)
			os.Exit(1)
		}
	}

	if *debug {
		spew.Fdump(os.Stderr, DLNADevice)
		clients, err := av1.NewAVTransport1ClientsByURL(DLNADevice.Location)
		spew.Fdump(os.Stderr, clients, err)
		for _, client := range clients {
			resp, err := http.Get(client.Location.String())
			if err != nil {
				spew.Fprintln(os.Stderr, err)
				continue
			}
			data, err := io.ReadAll(resp.Body)
			if err != nil {
				spew.Fprintln(os.Stderr, err)
				continue
			}
			spew.Fprintln(os.Stderr, string(data))
		}
		if !*headers {
			os.Exit(0)
		}
	}
	if *device == "" {
		fmt.Println("----------")
	}

	sink, err := chooseAudioSource(*source)
	if err != nil {
		fmt.Fprintln(os.Stderr, "audio:", err)
		os.Exit(1)
	}
	// on-demand handling of blast sink
	if sink == BLASTMONITOR {
		blastSink := exec.Command(
			"pactl", "load-module", "module-null-sink", "sink_name=blast",
		)
		var err error
		blastSinkID, err = blastSink.Output()
		if err != nil {
			fmt.Fprintln(os.Stderr, "blast sink:", err)
			os.Exit(1)
		}
		blastSinkID = bytes.TrimSpace(blastSinkID)
	}

	if *source == "" {
		fmt.Println("----------")
	}
	streamHost, err := chooseStreamIP(*ip)
	if err != nil {
		fmt.Fprintln(os.Stderr, "network:", err)
		cleanup()
		os.Exit(1)
	}
	if *ip == "" {
		fmt.Println("----------")
	}

	log.Printf(
		"starting the stream on port %d "+
			"(configure your firewall if necessary)",
		*port,
	)
	streamHandler := stream{
		sink:         sink,
		mime:         *mime,
		format:       *format,
		bitrate:      *bitrate,
		chunk:        *chunk,
		printheaders: *headers,
		bitdepth:     *bits,
		samplerate:   *rate,
		channels:     *channels,
		nochunked:    *nochunked,
		bige:         *bige,
	}
	if *useaac {
		streamHandler.format = "adts"
		streamHandler.mime = "audio/aac"
	}
	if *useflac {
		streamHandler.format = "flac"
		streamHandler.mime = "audio/flac"
		streamHandler.bitrate = 0
	}
	if *uselpcm {
		streamHandler.format = "lpcm"
		streamHandler.mime = fmt.Sprintf("audio/L%d;rate=%d;channels=%d", *bits, *rate, *channels)
		streamHandler.bitrate = 0
	}
	if *usewav {
		streamHandler.format = "wav"
		streamHandler.mime = "audio/wav"
		streamHandler.bitrate = 0
	}

	streamHandler.contentfeat = dlnaContentFeatures{
		profileName:     strings.ToUpper(streamHandler.format),
		supportTimeSeek: true,
		supportRange:    false,
		flags: DLNA_ORG_FLAG_DLNA_V15 |
			DLNA_ORG_FLAG_CONNECTION_STALL |
			DLNA_ORG_FLAG_STREAMING_TRANSFER_MODE |
			DLNA_ORG_FLAG_BACKGROUND_TRANSFERT_MODE,
	}

	streamPath := "stream." + strings.ToLower(streamHandler.format)

	mux := http.NewServeMux()
	mux.Handle("/"+streamPath, streamHandler)
	var logoHandler logo = logobytes
	mux.Handle("/"+LOGO_PATH, logoHandler)
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", *port),
		ReadTimeout:  -1,
		WriteTimeout: -1,
		Handler:      mux,
	}
	go func() {
		err := httpServer.ListenAndServe()
		if err != nil {
			fmt.Fprintln(os.Stderr, "server:", err)
			cleanup()
			os.Exit(1)
		}
	}()
	// detect when the stream server is up
	for {
		_, err := net.Dial("tcp", fmt.Sprintf(":%d", *port))
		if err == nil {
			break
		}
	}

	var (
		streamURI string
		logoURI   string
		protocol  = "http"
	)

	if !*dummy && *format == "mp3" && detectSonos(DLNADevice) {
		protocol = "x-rincon-mp3radio"
	}

	if streamHost.To4() != nil {
		streamURI = fmt.Sprintf("%s://%s:%d/%s",
			protocol, streamHost, *port, streamPath)
		logoURI = fmt.Sprintf("http://%s:%d/%s",
			streamHost, *port, LOGO_PATH)
	} else {
		var zone string
		if streamHost.IsLinkLocalUnicast() {
			ifname, err := findInterface(streamHost)
			if err == nil {
				zone = "%" + ifname
			}
		}
		streamURI = fmt.Sprintf("%s://[%s%s]:%d/%s",
			protocol, streamHost, zone, *port, streamPath)
		logoURI = fmt.Sprintf("http://[%s%s]:%d/%s",
			streamHost, zone, *port, LOGO_PATH)
	}

	log.Printf("stream URI: %s\n", streamURI)

	log.Println("setting av1transport URI and playing")
	if !*dummy {
		av := av1setup{
			location:  DLNADevice.Location,
			stream:    streamHandler,
			logoURI:   logoURI,
			streamURI: streamURI,
		}
		err = AV1SetAndPlay(av)
		if err != nil {
			fmt.Fprintln(os.Stderr, "transport:", err)
			cleanup()
			os.Exit(1)
		}
	}

	isPlaying = true
	select {}
}
