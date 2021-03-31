package main

import "net"

type User struct {
	Name    string
	Address string
	C       chan string
	Conn    net.Conn
}

// create user
func NewUser(conn net.Conn) *User {
	userAddress := conn.RemoteAddr().String()

	user := &User{
		Name:    "编号:" + userAddress,
		Address: userAddress,
		C:       make(chan string),
		Conn:    conn,
	}

	// 启动协程去监听管道有无消息传入
	go user.ListenMsg()

	return user
}

// listen to user channel
// if msg was in, the msg will be sent to the target client
func (u *User) ListenMsg() {
	for {
		msg := <-u.C
		u.Conn.Write([]byte(msg + "\n"))
	}
}
