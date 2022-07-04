package main

import (
	"bufio"
	"fmt"
	"io"
	"net"

	"github.com/gethinyan/go-nat-traversal/pkg"
)

func main() {
	// 跟控制服务建立连接
	conn, err := net.Dial("tcp", pkg.RemoteAddr+":"+pkg.ControlServerPort)
	if err != nil {
		fmt.Println("「连接失败」")
		return
	}
	fmt.Println("「连接成功」:" + conn.RemoteAddr().String())

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil || err == io.EOF {
			break
		}
		if message != pkg.Connection {
			continue
		}
		// 跟本地服务建立连接
		connLeft := localConn()
		// 跟隧道服务建立连接
		connRigth := tunnelConn()
		// 转发连接
		pkg.Forward(connLeft, connRigth)
	}
	fmt.Println("「断开连接」")
}

func localConn() net.Conn {
	conn, err := net.Dial("tcp", ":"+pkg.LocalPort)
	if err != nil {
		fmt.Println("「连接本地服务失败」")
		return nil
	}
	fmt.Println("「连接本地服务成功」")

	return conn
}

func tunnelConn() net.Conn {
	conn, err := net.Dial("tcp", pkg.RemoteAddr+":"+pkg.TunnelServerPort)
	if err != nil {
		fmt.Println("「连接隧道失败」")
		return nil
	}
	fmt.Println("「连接隧道成功」")

	return conn
}
