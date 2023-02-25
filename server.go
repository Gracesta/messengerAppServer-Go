package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip         string
	Port       int
	serverChan chan string

	// mapLock and online user map
	OnlineUserMap map[string]*User
	mapLock       sync.RWMutex
}

// Server interface
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:            ip,
		Port:          port,
		OnlineUserMap: make(map[string]*User),
		serverChan:    make(chan string),
	}

	return server
}

func (server *Server) BroadCast(user *User, msg string) {
	// Broadcasting new online user
	announce := "[" + user.UserAddr + "]" + user.UserName + ":" + msg

	server.serverChan <- announce

}

func (server *Server) ListenServerChannel() {
	// Once new message got in server channel, anounce it to all online users

	for {
		msg := <-server.serverChan

		server.mapLock.Lock()
		for _, userCli := range server.OnlineUserMap {
			userCli.userChan <- msg
		}
		server.mapLock.Unlock()

	}

}

func (server *Server) Handler(conn net.Conn) {
	fmt.Println("Connection built:", conn)

	// Once a connection built, assign a user class for this
	user := NewUser(conn)

	// New user added to online user map
	server.mapLock.Lock()
	server.OnlineUserMap[user.UserName] = user
	server.mapLock.Unlock()

	// Online anouncement for new user
	server.BroadCast(user, "ONLINE")

	// select{}
}

func (server *Server) Start() {
	// socket listening
	Listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.Ip, server.Port))
	if err != nil {
		fmt.Println("socket Listen error:", err)
		return
	}
	fmt.Println("Server launched")
	defer Listener.Close()

	for {
		conn, err := Listener.Accept()
		if err != nil {
			fmt.Println("Listener accept error", err)
			continue
		}

		// parrallel
		go server.ListenServerChannel()

		go server.Handler(conn)
	}

}
