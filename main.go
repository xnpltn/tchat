package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)
const PORT = 3033

type message struct {
	text string
	sender net.Conn
}

type server struct {
	message chan message
	port    uint16
	conns map[string]net.Conn
}



func newServer(port uint16) *server {
	return &server{
		message: make(chan message),
		port:    port,
		conns: make(map[string]net.Conn),
	}
}


func (s *server) handleConns(conn net.Conn, msg chan message) {
	s.conns[conn.RemoteAddr().String()] = conn
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if errors.Is(err, io.EOF) {
				continue
			}
			log.Fatalf("error: %v", err)
		}
		msg <- message{
			text: string(buffer[:n]),
			sender: conn,
		}
	}

}

func (s *server) handlemessages(msg chan message) {
	
	for {
		message := <-msg
		fmt.Println(message.sender.RemoteAddr())
		fmt.Println(len(s.conns))
	}
}

func main() {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("localhost:%d", PORT))
	if err!= nil{
		fmt.Println("error: ", err)
		os.Exit(1)
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	server := newServer(PORT)
	message := make(chan message)
	go server.handlemessages(message)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
		}
		go server.handleConns(conn, message)

	}
}
