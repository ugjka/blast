package main

import "fmt"

const DLNA_ORG_FLAG_SENDER_PACED = (1 << 31)
const DLNA_ORG_FLAG_TIME_BASED_SEEK = (1 << 30)
const DLNA_ORG_FLAG_BYTE_BASED_SEEK = (1 << 29)
const DLNA_ORG_FLAG_PLAY_CONTAINER = (1 << 28)
const DLNA_ORG_FLAG_S0_INCREASE = (1 << 27)
const DLNA_ORG_FLAG_SN_INCREASE = (1 << 26)
const DLNA_ORG_FLAG_RTSP_PAUSE = (1 << 25)
const DLNA_ORG_FLAG_STREAMING_TRANSFER_MODE = (1 << 24)
const DLNA_ORG_FLAG_INTERACTIVE_TRANSFERT_MODE = (1 << 23)
const DLNA_ORG_FLAG_BACKGROUND_TRANSFERT_MODE = (1 << 22)
const DLNA_ORG_FLAG_CONNECTION_STALL = (1 << 21)
const DLNA_ORG_FLAG_DLNA_V15 = (1 << 20)
const DLNA_ORG_FLAG_LINK_PROTECTED = (1 << 16)
const DLNA_ORG_FLAG_CLEAR_TEXT_BYTE_SEEK_FULL = (1 << 15)
const DLNA_ORG_FLAG_CLEAR_TEXT_BYTE_SEEK_LIMITED = (1 << 14)

func formatDLNAFlags(flags int) string {
	return fmt.Sprintf("DLNA.ORG_FLAGS=%.8x%.24x", flags, 0)
}

type dlnaContentFeatures struct {
	profileName     string
	supportTimeSeek bool
	supportRange    bool
	transcoded      bool
	flags           int
}

func (c dlnaContentFeatures) String() (out string) {
	if c.profileName != "" {
		out += fmt.Sprintf("DLNA.ORG_PN=%s;", c.profileName)
	}
	if c.supportTimeSeek || c.supportRange {
		out += fmt.Sprintf("DLNA.ORG_OP=%d%d;", bti(c.supportTimeSeek), bti(c.supportRange))
	}
	if c.transcoded {
		out += fmt.Sprintf("DLNA.ORG_CI=%d;", bti(c.transcoded))
	}
	out += formatDLNAFlags(c.flags)
	return
}

func bti(b bool) int {
	if b {
		return 1
	}
	return 0
}
