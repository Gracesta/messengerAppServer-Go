package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

// Server interface
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
	}

	return server
}

func (server *Server) Handler(conn net.Conn) {
	fmt.Println("Connection built:", conn)
}

func (server *Server) Start() {
	// socket listening
	Listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.Ip, server.Port))
	if err != nil {
		fmt.Println("socket Listen error:", err)
		return
	}

	defer Listener.Close()
	for {
		conn, err := Listener.Accept()
		if err != nil {
			fmt.Println("Listener accept error", err)
			continue
		}

		// handler parrallel
		go server.Handler(conn)
	}

}
