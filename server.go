package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip string
	Port int
}

// NewServer creates a new server instance
func NewServer(ip string, port int) *Server {
	return &Server{
		Ip: ip,
		Port: port,
	}
}

func (this *Server) Handler(conn net.Conn)  {
	// logic
	fmt.Println("connect success")
}

// Start starts the server
func (this *Server) Start()  {
	// socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	// close listen socket
	defer listener.Close()

	for {
		// accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept err:", err)
			continue
		}

		// do handler
		go this.Handler(conn)
	}
}