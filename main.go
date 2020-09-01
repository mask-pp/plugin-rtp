package rtp

import (
	"fmt"
	. "github.com/Monibuca/engine"
	. "github.com/mask-pp/plugin-rtp/transport"
	"golang.org/x/sync/errgroup"
	"log"
)

var (
	config = struct {
		ListenPorts string
		AutoPull    bool
		Timeout     int
		Reconnect   bool
	}{"6000:6000", false, 0, false}
)

func init() {
	InstallPlugin(&PluginConfig{
		Name:   "RTP",
		Type:   PLUGIN_PUBLISHER | PLUGIN_HOOK,
		Config: &config,
		Run:    runPlugin,
		HotConfig: map[string]func(interface{}){
			"AutoPull": func(value interface{}) {
				config.AutoPull = value.(bool)
			},
		},
	})
}

func runPlugin() {
	start, end, err := SplitPorts(config.ListenPorts)
	if err != nil {
		log.Fatal(err)
	}

	var e errgroup.Group
	for ; start <= end; start++ {
		start := start
		e.Go(func() error {
			if IsPortInUse(start) {
				return fmt.Errorf("RTSP port[%d] In Use", start)
			}
			return ListenRtp(uint16(start))
		})
	}

	if err = e.Wait(); err != nil {
		log.Fatal(err)
	}
}

// 启动一套rtp数据接收生产线
func ListenRtp(port uint16) error {
	var (
		e   errgroup.Group
		rtp = &Session{
			ITransport: NewTCPServer(port, true),
			RTP:        RTP{},
			PSCache:    make(map[uint32]*PSInfo),
		}
	)
	defer rtp.Close()

	e.Go(func() error {
		return rtp.Start()
	})

	e.Go(func() error {
		return rtp.PacketHandler()
	})

	return e.Wait()
}
