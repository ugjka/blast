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
	"encoding/json"
	"fmt"
	"os/exec"
)

func chooseAudioSource(lookup string) (string, error) {
	srcCMD := exec.Command("pactl", "-f", "json", "list", "sources", "short")
	srcData, err := srcCMD.Output()
	if err != nil {
		return "", fmt.Errorf("pactl sources: %v", err)
	}

	var srcJSON Sources
	err = json.Unmarshal(srcData, &srcJSON)
	if err != nil {
		return "", err
	}
	if len(srcJSON) == 0 {
		return "", fmt.Errorf("no audio sources found")
	}
	srcJSON = append(srcJSON, struct{ Name string }{BLASTMONITOR})
	if lookup != "" {
		for _, v := range srcJSON {
			if v.Name == lookup {
				return lookup, nil
			}
		}
		return "", fmt.Errorf("%s: not found", lookup)
	}

	fmt.Println("Audio sources")
	// append for on-demand loading of blast sink
	for i, v := range srcJSON {
		fmt.Printf("%d: %s\n", i, v.Name)
	}

	fmt.Println("----------")
	fmt.Println("Select the audio source:")

	selected := selector(srcJSON)
	return srcJSON[selected].Name, nil
}

type Sources []struct {
	Name string
}
