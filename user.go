package main

import (
	"net"
)

type User struct {
	Name    string
	Address string
	UserC   chan string
	Conn    net.Conn
	Server  *Serve
}

// create user
func NewUser(conn net.Conn, server *Serve) *User {
	userAddress := conn.RemoteAddr().String()

	user := &User{
		Name:    userAddress,
		Address: userAddress,
		UserC:   make(chan string),
		Conn:    conn,
		Server:  server,
	}

	// 启动协程去监听管道有无消息传入
	go user.ListenMsg()

	return user
}

// 用户上线
func (u *User) Online() {
	u.Server.MapLock.Lock()
	u.Server.OnlineUserMap[u.Name] = u
	u.Server.MapLock.Unlock()

	// 广播用户上线消息
	u.Server.BroadCast(u, "上线")
}

// 用户下线
func (u *User) Offline() {
	u.Server.MapLock.Lock()
	delete(u.Server.OnlineUserMap, u.Name)
	u.Server.MapLock.Unlock()

	// 广播用户下线消息
	u.Server.BroadCast(u, "下线")
}

// 用户消息传递
func (u *User) sendMsg(msg string) {
	u.Server.BroadCast(u, msg)
}

// listen to user channel
// if msg was in, the msg will be sent to the target client
func (u *User) ListenMsg() {
	for {
		msg := <-u.UserC
		u.Conn.Write([]byte(msg + "\n"))
	}
}
