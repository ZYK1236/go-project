package main

import (
	"fmt"
	"net"
)

type Serve struct {
	Ip   string
	Port int
}

// 创建 server 的方法
func NewServer(ip string, port int) *Serve {
	server := &Serve{
		Ip:   ip,
		Port: port,
	}

	return server
}

// conn handler 逻辑代码
func (s *Serve) Handle(conn net.Conn) {
	fmt.Println("listen success")
}

// Serve 的启动方法
func (s *Serve) Start() {
	// socket start
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("net listen err:", err)
		return
	}

	// accept data
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("conn err:", err)
		}

		// handler，开辟一个协程，不堵塞
		go s.Handle(conn)
	}
}
