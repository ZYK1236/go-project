package main

import (
	"flag"
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	Conn       net.Conn
}

func NewClient(serverIp string, serverPort int) *Client {
	// 创建客户端
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
	}
	// 客户端发送请求给服务端要求建立链接
	conn, error := net.Dial("tcp", fmt.Sprintf("%s:%d", client.ServerIp, client.ServerPort))
	if error != nil {
		fmt.Println("net.Dial error:", error.Error())
		return nil
	}
	client.Conn = conn

	return client
}

var serverIp string
var serverPort int

// init 在 main 之前执行
func init() {
	// 做一个数据与命令行的绑定
	// 例如输入 -ip 127.0.0.1 -port 8888
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "默认地址")
	flag.IntVar(&serverPort, "port", 8888, "默认端口")
}

func main() {
	// 命令行解析
	flag.Parse()
	fmt.Println(serverIp, serverPort)
	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>>>> NewClient error <<<<<<<<")
		return
	}
	fmt.Println("NewClient success")

	select {}
}
