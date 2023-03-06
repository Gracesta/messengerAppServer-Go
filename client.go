package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
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
			client.PublicChat()
			break
		case 2:
			client.PrivateChat()
			break
		case 3:
			client.UpdateName()
			break
		}
	}
}

func (client *Client) SelectUser() {

	// send who to know the user online
	_, err := client.conn.Write([]byte("who\n"))
	if err != nil {
		fmt.Println("client conn.write error: ", err)
		return
	}
}

func (client *Client) PrivateChat() {

	// show online user first
	client.SelectUser()

	var targetUserName string
	var chatMsg string

	fmt.Println(">>>>> Input the name of user who you wanna chat with")
	fmt.Scanln(&targetUserName)

	for targetUserName != "exit" {
		var errBufio error

		fmt.Println(">>>>> Input the message you wanna send (input \"exit\" to exit)")
		in := bufio.NewReader(os.Stdin)
		chatMsg, errBufio = in.ReadString('\n')
		if errBufio != nil {
			fmt.Println("reading string error", errBufio)
		}
		// fmt.Scanln(&chatMsg)

		// send message to usr conn to broadcast
		for chatMsg != "exit\r\n" {
			if len(chatMsg) != 0 {
				sendMsg := "to|" + targetUserName + "|" + chatMsg + "\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("client conn.write error: ", err)
					break
				}

				chatMsg = ""
				fmt.Println(">>>>> Input the message you wanna send (input \"exit\" to exit)")
				chatMsg, errBufio = in.ReadString('\n')
				if errBufio != nil {
					fmt.Println("reading string error", errBufio)
				}
				// fmt.Scanln(&chatMsg)
			}
		}

		client.SelectUser()
		fmt.Println(">>>>> Input the name of user who you want chat with (input \"exit\" to exit)")
		fmt.Scanln(&targetUserName)
	}
}

func (client *Client) PublicChat() {
	var chatMsg string
	var errBufio error
	in := bufio.NewReader(os.Stdin)

	fmt.Println(">>>>> Input the message you want to send (input \"exit\" to exit)")
	chatMsg, errBufio = in.ReadString('\n')
	if errBufio != nil {
		fmt.Println("reading string error", errBufio)
	}
	// fmt.Scanln(&chatMsg)

	for chatMsg != "exit\r\n" {

		// send message to usre conn to broadcast
		if len(chatMsg) != 0 {
			sendMsg := chatMsg
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("client conn.write error: ", err)
				break
			}

		}

		chatMsg = ""
		fmt.Println(">>>>> Input the message you want to send (input \"exit\" to exit)")
		chatMsg, errBufio = in.ReadString('\n')
		if errBufio != nil {
			fmt.Println("reading string error", errBufio)
		}
		// fmt.Scanln(&chatMsg)
	}
}

func (client *Client) UpdateName() bool {
	fmt.Println("Input the new username you want")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("client conn.write error: ", err)
		return false
	}

	return true

}

func (client *Client) HandleResponse() {
	// handle the message from server

	// Once read message from conn, display it
	io.Copy(os.Stdout, client.conn)

}

func (client *Client) menu() bool {
	var flag int

	fmt.Println("1. Public Chat")
	fmt.Println("2. Private Chat")
	fmt.Println("3. Update Username")
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

	// go routin for handling message from server
	go client.HandleResponse()

	client.Run()
}
