package main

import "net"

type User struct {
	UserName string
	UserAddr string
	userChan chan string
	conn     net.Conn // every user has an unique connection to their client

	server *Server
}

func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		UserName: userAddr,
		UserAddr: userAddr,
		userChan: make(chan string),
		conn:     conn,
		server:   server,
	}

	// Once new user built, set a go routine for it to listen message
	go user.ListenUserMessage()

	return user
}

// Listen user's channel, send message once received
func (user *User) ListenUserMessage() {
	for {
		msg := <-user.userChan
		user.conn.Write([]byte(msg + "\n"))
	}
}

func (user *User) Online() {
	// New user added to online user map
	user.server.mapLock.Lock()
	user.server.OnlineUserMap[user.UserName] = user
	user.server.mapLock.Unlock()

	// Online anouncement for new user
	user.server.BroadCast(user, "ONLINE")

}

func (user *User) Offline() {
	// New user added to online user map
	user.server.mapLock.Lock()
	delete(user.server.OnlineUserMap, user.UserName)
	user.server.mapLock.Unlock()

	// Online anouncement for new user
	user.server.BroadCast(user, "OFFLINE")
}

func (user *User) DoMessage(msg string) {
	user.server.BroadCast(user, msg)
}
