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
	"net/url"
	"time"

	"github.com/huin/goupnp/dcps/av1"
)

func AV1SetAndPlay(loc *url.URL, stream string) {
	client, err := av1.NewAVTransport1ClientsByURL(loc)
	stderr(err)
	err = client[0].SetAVTransportURI(0, stream, "")
	stderr(fmt.Errorf("set uri: %v", err))
	time.Sleep(time.Second)
	err = client[0].Play(0, "1")
	stderr(fmt.Errorf("play: %v", err))
}

func AV1Stop(loc *url.URL) {
	client, err := av1.NewAVTransport1ClientsByURL(loc)
	if err != nil {
		return
	}
	client[0].Stop(0)
}
