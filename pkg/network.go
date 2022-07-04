package pkg

import (
	"io"
	"net"
)

// 常量定义
const (
	RemoteAddr        = "127.0.0.1" // "154.222.29.27"
	ControlServerPort = "9090"
	TCPServerPort     = "9091"
	TunnelServerPort  = "9092"
	LocalPort         = "80"
	Ping              = "ping\n"
	Connection        = "connection\n"
)

// Forward 转发
func Forward(connLeft, connRight net.Conn) {
	defer connLeft.Close()
	defer connRight.Close()

	go func() {
		io.Copy(connLeft, connRight)
	}()
	io.Copy(connRight, connLeft)
}
