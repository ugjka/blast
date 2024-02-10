package main

import (
	"testing"
)

func TestDidlMetadata(t *testing.T) {
	want := "&lt;DIDL-Lite xmlns=&#34;urn:schemas-upnp-org:metadata-1-0/DIDL-Lite/&#34; xmlns:dc=&#34;http://purl.org/dc/elements/1.1/&#34; xmlns:sec=&#34;http://www.sec.co.kr/dlna&#34; xmlns:upnp=&#34;urn:schemas-upnp-org:metadata-1-0/upnp/&#34;&gt; &lt;item id=&#34;0&#34; parentID=&#34;-1&#34; restricted=&#34;false&#34;&gt; &lt;res protocolInfo=&#34;http-get:*:audio/mpeg:*&#34;&gt;http://192.168.1.225:9000/stream.mp3&lt;/res&gt; &lt;/item&gt; &lt;/DIDL-Lite&gt;"
	got := didlMetadata("http://192.168.1.225:9000/stream.mp3")
	if want != got {
		t.Fatalf("\nwant: %s\n_got:%s", want, got)
	}
}
