package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

type message struct {
	text   string
	sender net.Conn
}

type server struct {
	message chan message
	port    uint16
	conns   map[string]net.Conn
}

func newServer(port uint16) *server {
	return &server{
		message: make(chan message),
		port:    port,
		conns:   make(map[string]net.Conn),
	}
}

func (s *server) handleConns(conn net.Conn, msg chan message) {
	s.conns[conn.RemoteAddr().String()] = conn
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if errors.Is(err, io.EOF) {
				delete(s.conns, conn.RemoteAddr().String())
				continue
			}
			log.Fatal(err)
		}
		msg <- message{
			text:   string(buffer[:n]),
			sender: conn,
		}
	}
}

func (s *server) handlemessages(msg chan message) {
	for {
		message := <-msg
		if len(message.text) > 1 {
			for _, conn := range s.conns {
				if conn.RemoteAddr().String() != message.sender.RemoteAddr().String() {
					_, err := conn.Write([]byte(message.text))
					if err != nil {
						log.Println(err)
					}
				}
			}
		} else {
			_, err := message.sender.Write([]byte("message can't be empy\n"))
			if err != nil {
				log.Println("error: ", err)
			}
		}
	}
}

func main() {
	port := flag.Int64("port", 6969, ` port is a number specifying the port which the application will run on`)
	flag.Parse()
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Println("error: ", err)
		os.Exit(1)
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Servr starting at : %s \n", listener.Addr().String())
	server := newServer(uint16(*port))
	messageChan := make(chan message)
	go server.handlemessages(messageChan)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		go server.handleConns(conn, messageChan)

	}
}
