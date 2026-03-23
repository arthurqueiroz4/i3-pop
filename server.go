package main

import (
	"errors"
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
			if errors.Is(err, net.ErrClosed) {
				return
			}
			if netErr, ok := err.(net.Error); ok && netErr.Temporary() {
				log.Printf("Temporary accept error: %v", err)
				continue
			}
			log.Printf("Failed to accept connection: %v", err)
			continue
		}
		buf := make([]byte, 4096)
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("Error reading from connection: %v", err)
			conn.Close()
			continue
		}

		msg := strings.TrimSpace(string(buf[:n]))
		go handler(msg)
		_, err = conn.Write([]byte("Message received and processing"))
		if err != nil {
			log.Printf("Error writing to connection: %v", err)
		}
		conn.Close()
	}
}

func (s *Server) Stop() {
	if s.listener != nil {
		s.listener.Close()
	}
}
