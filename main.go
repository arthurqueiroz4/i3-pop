package main

import (
	"log"
	"net"
	"strings"
	"sync"

	"github.com/golang-collections/collections/stack"
	"go.i3wm.org/i3/v4"
)

func main() {
	var (
		s     = stack.New()
		mutex sync.Mutex
	)

	go func() {
		subscription := i3.Subscribe(i3.WorkspaceEventType)
		defer subscription.Close()

		for subscription.Next() {
			event := subscription.Event().(*i3.WorkspaceEvent)
			if event.Change == "focus" {
				mutex.Lock()
				if event.Old.Name != "" {
					s.Push(event.Old.Name)
					log.Printf("[ %s ] - Workspace pushed to stack", event.Old.Name)
				}
				mutex.Unlock()
			}
		}

		if err := subscription.Close(); err != nil {
			log.Printf("Error closing subscription: %v", err)
		}
	}()

	listener, err := net.Listen("tcp", ":43222")
	if err != nil {
		log.Fatalf("Failed to start TCP listener: %v", err)
	}
	defer listener.Close()
	log.Println("Socket listening on port 43222")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}
		log.Println("New connection established")

		go handleConnection(conn, s, &mutex)
	}
}

func handleConnection(conn net.Conn, s *stack.Stack, mutex *sync.Mutex) {
	defer conn.Close()

	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		log.Printf("Error reading from connection: %v", err)
		return
	}

	msg := strings.TrimSpace(string(buf[:n]))
	log.Printf("Message received: %s", msg)

	if msg == "back" {
		mutex.Lock()
		defer mutex.Unlock()

		if s.Len() == 0 {
			log.Println("Stack is empty, no workspace to go back to")
			return
		}

		workspaceName := s.Pop().(string)
		log.Printf("Switching to workspace: %s", workspaceName)
		if _, err := i3.RunCommand("workspace " + workspaceName); err != nil {
			log.Printf("Failed to switch workspace: %v", err)
		}
	} else {
		log.Printf("Invalid command: %s", msg)
	}
}