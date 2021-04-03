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
	flag       int // 当前用户选择模式
}

func NewClient(serverIp string, serverPort int) *Client {
	// 创建客户端
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
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

func Menu(client *Client) bool {
	var flag int
	fmt.Println("1. 公共聊天模式")
	fmt.Println("2. 私人聊天模式")
	fmt.Println("3. 重命名")
	fmt.Println("0. 退出")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		client.flag = flag
		fmt.Println("你选择了:", flag)
		return true
	} else {
		fmt.Println("输入合法的选项")
		return false
	}
}

func (c *Client) Run() {
	for {
		if c.flag != 0 {
			for {
				isChoose := Menu(c)
				if isChoose {
					break
				}
			}
			// 根据不同的选择模式去处理
			// switch 默认情况下 case 最后自带 break 语句，匹配成功后就不会执行其他 case，如果我们需要执行后面的 case，可以使用 fallthrough
			switch c.flag {
			case 1:
				fmt.Println(">>>>>>>> choose 公聊模式 <<<<<<<<")
			case 2:
				fmt.Println(">>>>>>>> choose 私聊模式 <<<<<<<<")
			case 3:
				fmt.Println(">>>>>>>> choose 重新命名 <<<<<<<<")
			}
		} else {
			break
		}
	}
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
	client.Run()
}
