package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial err:", err)
		return nil
	}
	client.conn = conn
	return client
}

func (client *Client) DealResponse() {
	io.Copy(os.Stdout, client.conn)
}

func (client *Client) menu() bool {
	var flag int
	fmt.Println("1. public chat")
	fmt.Println("2. private chat")
	fmt.Println("3. rename")
	fmt.Println("0. exit")
	fmt.Scanln(&flag)
	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println(">>>>>> invalid input")
		return false
	}
}

func (client *Client) PublicChat() {
	var chatMsg string
	fmt.Println(">>>>>> input message")
	fmt.Scanln(&chatMsg)
	for chatMsg != "exit" {
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn.Write err:", err)
				break
			}
		}
		chatMsg = ""
		fmt.Println(">>>>>> input message")
		fmt.Scanln(&chatMsg)
	}
}

func (client *Client) SelectUsers() {
	sendMsg := "who\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return
	}
}

func (client *Client) PrivateChat() {
	var remoteName string
	var chatMsg string

	client.SelectUsers()
	fmt.Println(">>>>>> input remote user")
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Println(">>>>>> input message")
		fmt.Scanln(&chatMsg)
		for chatMsg != "exit" {
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteName + "|" + chatMsg + "\n\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn.Write err:", err)
					break
				}
			}
			chatMsg = ""
			fmt.Println(">>>>>> input message")
			fmt.Scanln(&chatMsg)
		}
		client.SelectUsers()
		fmt.Println(">>>>>> input remote user")
		fmt.Scanln(&remoteName)
	}
}

func (client *Client) Rename() bool {
	fmt.Println(">>>>>> input new name")
	fmt.Scanln(&client.Name)
	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return false
	}
	return true
}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {
		}
		switch client.flag {
		case 1:
			fmt.Println(">>>>>> public chat")
			client.PublicChat()
			break
		case 2:
			fmt.Println(">>>>>> private chat")
			client.PrivateChat()
			break
		case 3:
			fmt.Println(">>>>>> rename")
			client.Rename()
			break
		}
	}
}

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "server ip address")
	flag.IntVar(&serverPort, "port", 8888, "server port")
}

func main() {
	flag.Parse()
	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>>>> connect server failed")
		return
	}
	go client.DealResponse()
	fmt.Println(">>>>>> connect server success")

	client.Run()
}
