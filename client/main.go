package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
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
	fmt.Println("9. 退出")

	fmt.Scanln(&flag)
	fmt.Println("flag:", flag)

	if flag >= 1 && flag <= 3 {
		client.flag = flag
		fmt.Println("你选择了:", flag)
		return true
	} else if flag == 9 {
		client.flag = 9
		fmt.Println("退出")
		return true
	} else {
		fmt.Println(">>>>> 请输入合法的选项 <<<<<<")
		return false
	}
}

func (c *Client) PublicChat() {
	var chatMsg string

	fmt.Println(">>>>>> 输入消息:")
	fmt.Scanln(&chatMsg)

	for chatMsg != "exit" {
		if len(chatMsg) != 0 {
			sendMessage := chatMsg + "\n"
			_, err := c.Conn.Write([]byte(sendMessage))
			if err != nil {
				fmt.Println("c.Conn.Write error:", err.Error())
				break
			}
		}
		chatMsg = ""
		fmt.Println(">>>>>> 输入消息:")
		fmt.Scanln(&chatMsg)
	}
}

// update name
func (c *Client) UpdateName() bool {
	fmt.Println(">>>>>> 输入你的新用户名:")
	fmt.Scanln(&c.Name)

	sendMessage := "rename|" + c.Name + "\n"
	_, err := c.Conn.Write([]byte(sendMessage))
	if err != nil {
		fmt.Println("c.Conn.Write error:", err.Error())
		return false
	}

	return true
}

// 处理 server 返回的数据，并返回给用户
func (c *Client) DealResponseData() {
	// 一直等待 conn 有无消息写入，如果有，打印出来
	io.Copy(os.Stdout, c.Conn)
}

func (c *Client) Run() {
	for {
		if c.flag != 9 {
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
				c.PublicChat()
			case 2:
				fmt.Println(">>>>>>>> choose 私聊模式 <<<<<<<<")
			case 3:
				fmt.Println(">>>>>>>> choose 重新命名 <<<<<<<<")
				c.UpdateName()
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

	// 处理 server 返回的消息
	go client.DealResponseData()

	fmt.Println("NewClient success")
	client.Run()
}
