package main

import (
	"flag"
	"fmt"
	"net"
)

type Client struct {
	Ip   string
	Port int
	Name string
	flag int

	conn net.Conn
}

func NewClient(ip string, port int) *Client {
	client := &Client{
		Ip:   ip,
		Port: port,
		flag: 999,
	}

	// Link to server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))

	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	}
	client.conn = conn

	return client
}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {
			// display menu while input illegal
		}

		switch client.flag {
		case 1:
			fmt.Println("1. Public Chat")
			break
		case 2:
			fmt.Println("2. Private Chat")
			break
		case 3:
			fmt.Println("3. Change Username")
			break
		case 0:
			fmt.Println("0. Exit")
			break
		}
	}
}

func (client *Client) menu() bool {
	var flag int

	fmt.Println("1. Public Chat")
	fmt.Println("2. Private Chat")
	fmt.Println("3. Change Username")
	fmt.Println("0. Exit")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println("Input number within range")
		return false
	}

}

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "Server IP")
	flag.IntVar(&serverPort, "port", 8888, "Server Port")
}

func main() {
	// command line parse
	flag.Parse()

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>>> Link to Server Falied <<<<<")
		return
	}

	fmt.Println(">>>>> Link to Server Succedded <<<<<")

	if client.flag > 0 && client.flag < 4 {
		client.menu()
	} else {
		fmt.Println("Input number within range")
	}

	client.Run()
}
