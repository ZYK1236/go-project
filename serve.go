package main

import (
	"fmt"
	"io"
	"net"
	"runtime"
	"sync"
	"time"
)

type Serve struct {
	Ip            string
	Port          int
	OnlineUserMap map[string]*User // key: string -> 根据用户名 key 找到对应 user
	MapLock       sync.RWMutex     // 锁
	Message       chan string      // 负责消息广播的通道
}

// 创建 server 的方法
func NewServer(ip string, port int) *Serve {
	server := &Serve{
		Ip:            ip,
		Port:          port,
		OnlineUserMap: make(map[string]*User),
		Message:       make(chan string),
	}

	return server
}

// 将消息存进负责消息广播的通道
func (s *Serve) BroadCast(user *User, msg string) {
	sendMsg := user.Name + ":" + msg
	s.Message <- sendMsg
}

// 监听广播消息通道，一旦有消息广播给其他用户
func (s *Serve) ListenMessager() {
	for {
		msg := <-s.Message
		s.MapLock.Lock()
		for _, usr := range s.OnlineUserMap {
			usr.UserC <- msg
		}
		s.MapLock.Unlock()
	}
}

// conn handler 逻辑代码
func (s *Serve) Handle(conn net.Conn) {
	// 用户上线，添加进 OnlineUserMap
	user := NewUser(conn, s)
	user.Online()

	// 监控当前用户是否活跃的 channel
	isLive := make(chan bool)

	// 设定时器
	timer := time.NewTimer(60 * time.Second)

	// 接受客户端发的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			// 读消息，如果读的消息长度为空，则下线
			len, err := conn.Read(buf)
			if len == 0 {
				user.Offline()
				runtime.Goexit()
			}
			if err != nil && err != io.EOF {
				fmt.Println("conn.Read err:", err)
			}

			msg := string(buf[:len-1])
			user.sendMsg(msg)

			isLive <- true
		}
	}()

	// 阻塞 + 超时处理
	for {
		select {
		case <-isLive:
			// 重置定时器
			timer.Reset(60 * time.Second)

		case <-timer.C:
			// 超时，踢出
			user.sendMsg("因为超时被踢了，真可怜")

			// 链接释放
			conn.Write([]byte("你已下线...\n"))
			conn.Close()

			// 释放该用户资源
			close(user.UserC)

			runtime.Goexit()
		}
	}

}

// 启动 serve
func (s *Serve) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("net listen err:", err)
		return
	}

	// 启动监听广播 message 的协程
	go s.ListenMessager()

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
