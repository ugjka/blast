# BLAST

![Blast Logo](logo.png)

## Stream your Linux audio to DLNA receivers

You need `pactl`, `parec` and `lame` executables/dependencies on your system to run Blast.

If you have all that then you can launch `blast` and it looks like this when you run it:

```
[user@user blast]$ ./blast 
----------
DLNA receivers
0: Kitchen
1: Phone
2: Bedroom
3: Livingroom TV
----------
Select the DLNA device:
[1]
----------
Audio sources
0: alsa_output.pci-0000_00_1b.0.analog-stereo.monitor
1: alsa_input.pci-0000_00_1b.0.analog-stereo
2: bluez_output.D8_AA_59_95_96_B7.1.monitor
3: blast.monitor
----------
Select the audio source:
[2]
----------
Your LAN ip addresses
0: 192.168.1.14
1: 192.168.122.1
2: 2a04:ec00:b9ab:555:3c50:e6e8:8ea:211f
3: 2a04:ec00:b9ab:555:806d:800b:1138:8b1b
4: fe80::f4c2:c827:a865:35e5
----------
Select the lan IP address for the stream:
[0]
----------
2023/07/08 23:53:07 starting the stream on port 9000 (configure your firewall if necessary)
2023/07/08 23:53:07 seting av1transport URI and playing
```

## Building

You need the `go` and `go-tools` toolchain, also `git`

then execute:

```
git clone https://github.com/ugjka/blast
cd blast
go build
```

now you can run blast with:
```
[user@user blast]$ ./blast
```

## Bins

Prebuilt Linux binaries are available on the releases [page](https://github.com/ugjka/blast/releases)

## Why not use pulseaudio-dlna?

This is for pipewire-pulse users.

## License

```
MIT+NoAI License

Copyright (c) 2023 ugjka <ugjka@proton.me>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights/
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

This code may not be used to train artificial intelligence computer models
or retrieved by artificial intelligence software or hardware.
```