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
	"os"
	"strconv"
)

func selector(total int) int {
	var choice int
	for {
		var choiceStr string
		_, err := fmt.Fscanln(os.Stdin, &choiceStr)
		if err != nil {
			fmt.Print("\033[1A\033[K")
			continue
		}
		choice, err = strconv.Atoi(choiceStr)
		if err != nil || choice >= total {
			fmt.Print("\033[1A\033[K")
		} else {
			break
		}
	}
	return choice
}
