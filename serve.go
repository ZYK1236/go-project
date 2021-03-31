package main

import (
	"fmt"
	"net"
	"sync"
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
	sendMsg := user.Name + msg
	s.Message <- sendMsg
}

// 监听广播消息通道，一旦有消息广播给其他用户
func (s *Serve) ListenMessager() {
	for {
		msg := <-s.Message
		s.MapLock.Lock()
		for _, usr := range s.OnlineUserMap {
			usr.C <- msg
		}
		s.MapLock.Unlock()
	}
}

// conn handler 逻辑代码
func (s *Serve) Handle(conn net.Conn) {
	// 用户上线，添加进 OnlineUserMap
	user := NewUser(conn)
	s.MapLock.Lock()
	s.OnlineUserMap[user.Name] = user
	s.MapLock.Unlock()

	// 广播用户上线消息
	s.BroadCast(user, "上线")
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
