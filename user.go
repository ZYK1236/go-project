package main

import (
	"net"
	"strings"
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
	u.writeMsg("请注意，如果持续一分钟不发送信息，会被强制下线..")
}

// 用户下线
func (u *User) Offline() {
	u.Server.MapLock.Lock()
	delete(u.Server.OnlineUserMap, u.Name)
	u.Server.MapLock.Unlock()

	// 广播用户下线消息
	u.Server.BroadCast(u, "下线")
}

// 当前用户下输出信息
func (u *User) writeMsg(msg string) {
	u.Conn.Write([]byte(msg + "\n"))
}

// 用户消息传递
func (u *User) sendMsg(msg string) {
	// 查询当前在线用户
	if msg == "who" {
		u.writeMsg("当前上线用户有:")

		u.Server.MapLock.Lock()
		for _, usr := range u.Server.OnlineUserMap {
			onlineUsr := usr.Address + "_" + usr.Name
			u.writeMsg(onlineUsr)
		}
		u.Server.MapLock.Unlock()

		return
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		// 修改用户名，修改格式 example: rename|zyk
		newName := strings.Split(msg, "|")[1]
		_, ok := u.Server.OnlineUserMap[newName]

		if ok {
			u.writeMsg("有用户重名了")
		} else {
			u.Server.MapLock.Lock()
			delete(u.Server.OnlineUserMap, u.Name)
			u.Name = newName
			u.Server.OnlineUserMap[newName] = u
			u.Server.MapLock.Unlock()

			u.writeMsg("成功修改名字，新名字为:" + newName)
		}

		return
	} else if len(msg) > 4 && msg[:3] == "to|" {
		// 私聊功能，格式为: to|zyk|yourMessage
		msgArr := strings.Split(msg, "|")
		targetName := msgArr[1]

		if targetName == "" {
			u.writeMsg("格式不对")
			return
		}
		remoteUser, ok := u.Server.OnlineUserMap[targetName]
		if !ok {
			u.writeMsg("用户不存在")
			return
		}

		message := msgArr[2]
		if message == "" {
			u.writeMsg("无发送内容")
			return
		}

		remoteUser.writeMsg(u.Name + "对您说" + message)

		return
	}

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
