package ps

import (
	"bytes"
	"github.com/32bitkid/bitreader"
)

type RtpParsePacket struct {
	*DecPSPackage
}

func NewRtpParsePacket() *RtpParsePacket {
	return &RtpParsePacket{
		DecPSPackage: &DecPSPackage{
			rawData: make([]byte, MAXFrameLen),
			rawLen:  0,
		},
	}
}

func (rtp *RtpParsePacket) Clean() {
	rtp.clean()
}

// data包含 接受到完整一帧数据后，所有的payload, 解析出去后是一阵完整的raw数据
func (rtp *RtpParsePacket) Read(data []byte) ([]byte, error) {

	// add the MPEG Program end code
	data = append(data, 0x00, 0x00, 0x01, 0xb9)
	br := bitreader.NewReader(bytes.NewReader(data))

	if rtp.DecPSPackage != nil {
		return rtp.decPackHeader(br, data)
	}

	return nil, nil
}
