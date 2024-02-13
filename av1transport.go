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
	"net/url"
	"strings"
	"time"

	"github.com/huin/goupnp/dcps/av1"
)

func AV1SetAndPlay(loc *url.URL, albumart, stream string) error {
	client, err := av1.NewAVTransport1ClientsByURL(loc)
	if err != nil {
		return err
	}

	try := func(metadata string) error {
		err = client[0].SetAVTransportURI(0, stream, metadata)
		if err != nil {
			return fmt.Errorf("set uri: %v", err)
		}
		time.Sleep(time.Second)
		err = client[0].Play(0, "1")
		if err != nil {
			return fmt.Errorf("play: %v", err)
		}
		return nil
	}

	metadata := didlMetadata(albumart, stream)
	err = try(metadata)
	if err == nil {
		return nil
	}
	log.Println(err)
	log.Println("trying without metadata")
	return try("")
}

func AV1Stop(loc *url.URL) {
	client, err := av1.NewAVTransport1ClientsByURL(loc)
	if err != nil {
		return
	}
	client[0].Stop(0)
}

func didlMetadata(albumart, stream string) string {
	out := fmt.Sprintf(didlTemplate, albumart, stream)
	out = strings.ReplaceAll(out, "\n", " ")
	out = strings.ReplaceAll(out, "> <", "><")
	return out
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
<res protocolInfo="http-get:*:audio/mpeg:*">%s</res>
</item>
</DIDL-Lite>`
