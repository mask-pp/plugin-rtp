package rtp

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/mask-pp/plugin-rtp/ps"
	. "github.com/mask-pp/plugin-rtp/transport"
	"github.com/pion/rtp"
	"log"
)

type PSInfo struct {
	rtp.Header
	isFull bool
	buffer []byte
}

type Session struct {
	ITransport
	RTP
	PSCache map[uint32]*PSInfo
}

func (session *Session) PacketHandler() error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("packet handler panic: ", err)
			PrintStack()
		}
	}()

	var (
		rptOpt   = ps.NewRtpParsePacket()
		ch       = session.ReadPacketChan()
		cache    *PSInfo
		psBuffer []byte
	)

	//阻塞读取消息
	for {
		// 清除数据
		rptOpt.Clean()
		select {
		case <-session.Done():
			return nil
		case p := <-ch:
			var (
				rtpPack = new(rtp.Packet)
				ssrc    uint32
				isFull  bool
			)
			if err := rtpPack.Unmarshal(p.Data); err != nil {
				log.Println(err)
				continue
			}

			ssrc = rtpPack.SSRC
			if hexutil.Encode(rtpPack.Payload[:4]) == "0x000001ba" {
				if _, ok := session.PSCache[ssrc]; !ok {
					session.PSCache[ssrc] = &PSInfo{Header: rtpPack.Header, isFull: false, buffer: make([]byte, 0)}
				} else {
					isFull = true
				}
			}
			cache = session.PSCache[ssrc]

			if isFull {
				psBuffer = cache.buffer
				cache.buffer = p.Data[12:]
			} else {
				cache.buffer = append(cache.buffer, p.Data[12:]...)
				continue
			}

			// 提取音视频数据
			data, err := rptOpt.Read(psBuffer)
			if err != nil {
				fmt.Printf("------------------- payload: %s\n", hexutil.Encode(data))
				log.Println(err)
				continue
			}
			fmt.Printf("--------------- video data: %#x", data)

			var pack *RTPPack
			switch {
			case rptOpt.VideoStreamType != 0:
				pack = &RTPPack{
					Type:      RTP_TYPE_VIDEO,
					Timestamp: rtpPack.Timestamp,
					Payload:   data,
				}
			case rptOpt.AudioStreamType != 0:
				pack = &RTPPack{
					Type:      RTP_TYPE_AUDIO,
					Timestamp: rtpPack.Timestamp,
					Payload:   data,
				}
			default:
				continue
			}
			session.PushPack(pack)
		}
	}
}
