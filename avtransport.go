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
	"log"
	"strings"
	"time"

	"github.com/huin/goupnp"
	"github.com/huin/goupnp/dcps/av1"
)

type avsetup struct {
	device    *goupnp.MaybeRootDevice
	stream    stream
	logoURI   string
	streamURI string
}

type avtransport interface {
	SetAVTransportURI(InstanceID uint32, CurrentURI string, CurrentURIMetaData string) (err error)
	Play(InstanceID uint32, Speed string) (err error)
	Stop(InstanceID uint32) (err error)
}

func AVSetAndPlay(av avsetup) error {
	var client avtransport
	switch {
	case strings.HasSuffix(av.device.USN, "AVTransport:1"):
		clients, err := av1.NewAVTransport1ClientsByURL(av.device.Location)
		if err != nil {
			return err
		}
		client = avtransport(clients[0])
	case strings.HasSuffix(av.device.USN, "AVTransport:2"):
		clients, err := av1.NewAVTransport2ClientsByURL(av.device.Location)
		if err != nil {
			return err
		}
		client = avtransport(clients[0])
	default:
		return fmt.Errorf("error: no avtransport found")
	}
	var err error
	try := func(metadata string) error {
		err = client.SetAVTransportURI(0, av.streamURI, metadata)
		if err != nil {
			return fmt.Errorf("set uri: %v", err)
		}
		time.Sleep(time.Second)
		err = client.Play(0, "1")
		if err != nil {
			return fmt.Errorf("play: %v", err)
		}
		return nil
	}
	metadata := fmt.Sprintf(
		didlTemplate,
		av.logoURI,
		av.stream.mime,
		av.stream.contentfeat,
		av.stream.bitdepth,
		av.stream.samplerate,
		av.stream.channels,
		av.streamURI,
	)
	metadata = strings.ReplaceAll(metadata, "\n", " ")
	metadata = strings.ReplaceAll(metadata, "> <", "><")
	err = try(metadata)
	if err == nil {
		return nil
	}
	log.Println(err)
	log.Println("trying without metadata")
	return try("")
}

func AVStop(device *goupnp.MaybeRootDevice) {
	var client avtransport
	switch {
	case strings.HasSuffix(device.USN, "AVTransport:1"):
		clients, err := av1.NewAVTransport1ClientsByURL(device.Location)
		if err != nil {
			return
		}
		client = avtransport(clients[0])
	case strings.HasSuffix(device.USN, "AVTransport:2"):
		clients, err := av1.NewAVTransport2ClientsByURL(device.Location)
		if err != nil {
			return
		}
		client = avtransport(clients[0])
	default:
		return
	}
	client.Stop(0)
}

const didlTemplate = `<DIDL-Lite
xmlns="urn:schemas-upnp-org:metadata-1-0/DIDL-Lite/"
xmlns:upnp="urn:schemas-upnp-org:metadata-1-0/upnp/"
xmlns:dc="http://purl.org/dc/elements/1.1/"
xmlns:dlna="urn:schemas-dlna-org:metadata-1-0/"
xmlns:sec="http://www.sec.co.kr/"
xmlns:pv="http://www.pv.com/pvns/">
<item id="0" parentID="-1" restricted="1">
<upnp:class>object.item.audioItem.musicTrack</upnp:class>
<dc:title>Audio Cast</dc:title>
<dc:creator>Blast</dc:creator>
<upnp:artist>Blast</upnp:artist>
<upnp:albumArtURI>%s</upnp:albumArtURI>
<res protocolInfo="http-get:*:%s:%s"
bitsPerSample="%d"
sampleFrequency="%d"
nrAudioChannels="%d">%s</res>
</item>
</DIDL-Lite>`
