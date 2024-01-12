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
	"encoding/xml"
	"io"
	"net/http"
	"strings"

	"github.com/huin/goupnp"
	"github.com/huin/goupnp/dcps/av1"
)

func detectSonos(dev *goupnp.MaybeRootDevice) bool {
	var xmldata struct {
		Device struct {
			Manufacturer string `xml:"manufacturer"`
		} `xml:"device"`
	}

	clients, err := av1.NewAVTransport1ClientsByURL(dev.Location)
	if err != nil {
		return false
	}
	for _, client := range clients {
		resp, err := http.Get(client.Location.String())
		if err != nil {
			return false
		}
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return false
		}
		err = xml.Unmarshal(data, &xmldata)
		if err != nil {
			return false
		}
		man := xmldata.Device.Manufacturer
		man = strings.ToLower(man)
		if strings.Contains(man, "sonos") {
			return true
		}
	}
	return false
}
