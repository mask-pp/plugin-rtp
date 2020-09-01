package rtp

import (
	. "github.com/Monibuca/engine"
)

type RTPType int

const (
	RTP_TYPE_AUDIO RTPType = iota
	RTP_TYPE_VIDEO
	RTP_TYPE_AUDIOCONTROL
	RTP_TYPE_VIDEOCONTROL
)

type RTPPack struct {
	Type      RTPType
	Timestamp uint32
	Payload   []byte
}

func (rt RTPType) String() string {
	switch rt {
	case RTP_TYPE_AUDIO:
		return "audio"
	case RTP_TYPE_VIDEO:
		return "video"
	case RTP_TYPE_AUDIOCONTROL:
		return "audio control"
	case RTP_TYPE_VIDEOCONTROL:
		return "video control"
	}
	return "unknow"
}

type RTP struct {
	NALU
}

func (rtp *RTP) PushPack(pack *RTPPack) {
	switch pack.Type {
	case RTP_TYPE_AUDIO:
		payload := pack.Payload
		auHeaderLen := (int16(payload[0]) << 8) + int16(payload[1])
		auHeaderLen = auHeaderLen >> 3
		auHeaderCount := int(auHeaderLen / 2)
		var auLenArray []int
		for iIndex := 0; iIndex < int(auHeaderCount); iIndex++ {
			auHeaderInfo := (int16(payload[2+2*iIndex]) << 8) + int16(payload[2+2*iIndex+1])
			auLen := auHeaderInfo >> 3
			auLenArray = append(auLenArray, int(auLen))
		}
		startOffset := 2 + 2*auHeaderCount
		for _, auLen := range auLenArray {
			endOffset := startOffset + auLen
			addHead := []byte{0xAF, 0x01}
			rtp.PushAudio(0, append(addHead, payload[startOffset:endOffset]...))
			startOffset = startOffset + auLen
		}
	case RTP_TYPE_VIDEO:
		rtp.WriteNALU(pack.Timestamp, pack.Payload)
	}
}
