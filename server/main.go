package main

import (
	"fmt"
	"net"
	"time"

	"github.com/gethinyan/go-nat-traversal/pkg"
)

var controlConn net.Conn

var connPool []net.Conn

func main() {
	// 控制服务
	go controlServer()
	// tcp 服务
	go tcpServer()
	// 隧道服务
	tunnelServer()
}

func controlServer() {
	l, err := net.Listen("tcp", "127.0.0.1:"+pkg.ControlServerPort)
	if err != nil {
		fmt.Println("「监听control服务失败」")
		return
	}
	fmt.Println("「监听control服务成功」")
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		if controlConn != nil {
			conn.Close()
			return
		}
		fmt.Println("「control新连接」" + conn.RemoteAddr().String())
		controlConn = conn

		go keepAlive()
	}
}

func keepAlive() {
	go func() {
		for {
			controlConn.Write([]byte(pkg.Ping))
			time.Sleep(10 * time.Second)
		}
	}()
}

func tcpServer() {
	l, err := net.Listen("tcp", "127.0.0.1:"+pkg.TCPServerPort)
	if err != nil {
		fmt.Println("「监听tcp服务成功」")
		return
	}
	fmt.Println("「监听tcp服务成功」")
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("「tcp新连接」" + conn.RemoteAddr().String())
		go handleTCPServer(conn)
	}
}

func handleTCPServer(conn net.Conn) {
	connPool = append(connPool, conn)
	if controlConn == nil {
		fmt.Println("「无已连接的客户端」")
		return
	}
	controlConn.Write([]byte(pkg.Connection))
}

func tunnelServer() {
	l, err := net.Listen("tcp", "127.0.0.1:"+pkg.TunnelServerPort)
	if err != nil {
		fmt.Println("「监听隧道端口成功」")
		return
	}
	fmt.Println("「监听隧道端口成功」")
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("「tunnel新连接」" + conn.RemoteAddr().String())
		go handleTunnelServer(conn)
	}
}

func handleTunnelServer(conn net.Conn) {
	if len(connPool) <= 0 {
		fmt.Println("「连接池无有效连接」")
		return
	}
	tcpConn := connPool[0]
	connPool = connPool[1:]

	go pkg.Forward(tcpConn, conn)
}
