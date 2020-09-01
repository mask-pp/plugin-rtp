package rtp

import (
	"fmt"
	"net"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func IsPortInUse(port int) bool {
	if conn, err := net.DialTimeout("tcp", net.JoinHostPort("", fmt.Sprintf("%d", port)), 3*time.Second); err == nil {
		conn.Close()
		return true
	}
	return false
}

func SplitPorts(ports string) (int, int, error) {
	portLst := strings.Split(ports, ":")
	if len(portLst) != 2 {
		return 0, 0, fmt.Errorf("rtsp port type is not right")
	}
	start, err := strconv.Atoi(strings.Trim(portLst[0], " "))
	if err != nil {
		return 0, 0, err
	}
	end, err := strconv.Atoi(strings.Trim(portLst[1], " "))
	if err != nil {
		return 0, 0, err
	}
	if end < start || start <= 0 {
		return 0, 0, fmt.Errorf("start or end port is not right")
	}
	return start, end, err
}

func PrintStack() {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	fmt.Printf("==> %s\n", string(buf[:n]))
}
