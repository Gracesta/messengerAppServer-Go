package main

import "net"

type User struct {
	UserName string
	UserAddr string
	userChan chan string
	conn     net.Conn // every user has an unique connection to their client
}

func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		UserName: userAddr,
		UserAddr: userAddr,
		userChan: make(chan string),
		conn:     conn,
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
