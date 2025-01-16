package main

import (
	"log"
	"net"
	"strings"
)

type Server struct {
	network  string
	port     string
	listener net.Listener
}

func NewServer(network string, port string) *Server {
	return &Server{network: network, port: port}
}

func (s *Server) Start(handler func(msgToProcess string), ready chan struct{}) {
	listener, err := net.Listen(s.network, s.port)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	s.listener = listener
	log.Println("Server listening on", s.port)
	close(ready)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			return
		}
		buf := make([]byte, 4096)
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("Error reading from connection: %v", err)
			return
		}

		msg := strings.TrimSpace(string(buf[:n]))
		go handler(msg)
		_, err = conn.Write([]byte("Message received and processing"))
		conn.Close()
	}
}

func (s *Server) Stop() {
	if s.listener != nil {
		s.listener.Close()
	}
}
