package main

import (
	"testing"
)

func TestContentFeatures(t *testing.T) {
	f := dlnaContentFeatures{
		profileName:     "MP3",
		supportTimeSeek: true,
		supportRange:    false,
		transcoded:      true,
		flags: DLNA_ORG_FLAG_STREAMING_TRANSFER_MODE |
			DLNA_ORG_FLAG_BACKGROUND_TRANSFERT_MODE |
			DLNA_ORG_FLAG_CONNECTION_STALL |
			DLNA_ORG_FLAG_DLNA_V15,
	}
	want := "DLNA.ORG_PN=MP3;DLNA.ORG_OP=10;DLNA.ORG_CI=1;" +
		"DLNA.ORG_FLAGS=01700000000000000000000000000000"
	if f.String() != want {
		t.Fatalf("got %s, wanted %s", f, want)
	}
}
